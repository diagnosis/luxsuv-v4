package main

import (
	"fmt"
	"time"

	"github.com/diagnosis/luxsuv-v4/internal/auth"
	"github.com/diagnosis/luxsuv-v4/internal/config"
	"github.com/diagnosis/luxsuv-v4/internal/handlers"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/diagnosis/luxsuv-v4/internal/middleware"
	"github.com/diagnosis/luxsuv-v4/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	// Initialize logger
	log, err := logger.NewLogger("app.log")
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer log.Close()

	// Load configuration
	cfg, err := config.LoadConfig(log)
	if err != nil {
		log.Err("Failed to load config: " + err.Error())
		return
	}
	log.Info("Configuration loaded successfully")

	// Run migrations
	db, err := sqlx.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Err("Failed to open database for migrations: " + err.Error())
		return
	}
	defer db.Close()

	goose.SetLogger(&GooseLogger{log: log})
	if err := goose.Up(db.DB, "migrations"); err != nil {
		log.Err("Failed to apply migrations: " + err.Error())
		return
	}
	log.Info("Database migrations applied successfully")

	// Connect to database with connection pool settings
	db, err = sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Err("Failed to connect to database: " + err.Error())
		return
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Err("Failed to ping database: " + err.Error())
		return
	}
	log.Info("Successfully connected to database")

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(db)
	authService := auth.NewService(userRepo, cfg.JWTSecret, log)
	authHandler := handlers.NewAuthHandler(authService, log)
	authMiddleware := middleware.NewAuthMiddleware(authService, log)

	// Set up Echo server
	e := echo.New()

	// Configure Echo
	e.HideBanner = true
	e.HidePort = true

	// Global middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	e.Use(echomiddleware.Secure())

	// Configure rate limiter
	generalRateLimiterConfig := echomiddleware.RateLimiterConfig{
		Store: echomiddleware.NewRateLimiterMemoryStoreWithConfig(
			echomiddleware.RateLimiterMemoryStoreConfig{
				Rate:      5,                   // 5 requests per second
				Burst:     10,                  // Allow burst of 10 requests
				ExpiresIn: 3 * time.Minute,     // Clean up expired entries
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(429, map[string]string{"error": "too many requests"})
		},
	}

	// More restrictive rate limiter for authentication endpoints
	authRateLimiterConfig := echomiddleware.RateLimiterConfig{
		Store: echomiddleware.NewRateLimiterMemoryStoreWithConfig(
			echomiddleware.RateLimiterMemoryStoreConfig{
				Rate:      2,                   // 2 requests per second
				Burst:     5,                   // Allow burst of 5 requests
				ExpiresIn: 5 * time.Minute,     // Clean up expired entries
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(429, map[string]string{"error": "too many requests"})
		},
	}

	// Apply general rate limiting
	e.Use(echomiddleware.RateLimiterWithConfig(generalRateLimiterConfig))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Public routes with stricter rate limiting
	authGroup := e.Group("")
	authGroup.Use(echomiddleware.RateLimiterWithConfig(authRateLimiterConfig))
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	// Protected routes
	protectedGroup := e.Group("")
	protectedGroup.Use(authMiddleware.RequireAuth())
	protectedGroup.GET("/users/me", authHandler.GetCurrentUser)

	// Admin routes
	adminGroup := e.Group("/admin")
	adminGroup.Use(authMiddleware.RequireAuth())
	adminGroup.Use(authMiddleware.RequireAdmin())
	adminGroup.DELETE("/users/:id", authHandler.DeleteUser)

	log.Info("Starting server on port " + cfg.Port)
	log.Info("Available endpoints:")
	log.Info("  GET  /health")
	log.Info("  POST /register")
	log.Info("  POST /login")
	log.Info("  GET  /users/me (protected)")
	log.Info("  DELETE /admin/users/:id (admin only)")

	if err := e.Start(":" + cfg.Port); err != nil {
		log.Err("Failed to start server: " + err.Error())
	}
}

// GooseLogger adapts your logger to Goose's logger interface
type GooseLogger struct {
	log *logger.Logger
}

func (g *GooseLogger) Fatal(v ...interface{}) {
	g.log.Err(fmt.Sprint(v...))
	panic(v)
}

func (g *GooseLogger) Fatalf(format string, v ...interface{}) {
	g.log.Err(fmt.Sprintf(format, v...))
	panic(fmt.Sprintf(format, v...))
}

func (g *GooseLogger) Print(v ...interface{}) {
	g.log.Info(fmt.Sprint(v...))
}

func (g *GooseLogger) Println(v ...interface{}) {
	g.log.Info(fmt.Sprint(v...))
}

func (g *GooseLogger) Printf(format string, v ...interface{}) {
	g.log.Info(fmt.Sprintf(format, v...))
}