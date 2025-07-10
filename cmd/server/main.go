package main

import (
	"fmt"
	"time"
	"github.com/diagnosis/luxsuv-v4/internal/config"
	"github.com/diagnosis/luxsuv-v4/internal/data"
	"github.com/diagnosis/luxsuv-v4/internal/handlers"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/diagnosis/luxsuv-v4/internal/mw"
	"github.com/diagnosis/luxsuv-v4/internal/services"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// Initialize dependencies
	repo := data.NewRepository(db)
	authService := services.NewAuthService(repo, cfg.JWTSecret, log)
	authHandler := handlers.NewAuthHandler(authService, log)

	// Set up Echo server
	e := echo.New()
	
	// Configure Echo
	e.HideBanner = true
	e.HidePort = true

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())

	// Configure rate limiter with more restrictive limits
	generalRateLimiterConfig := middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      5,                    // 5 requests per second
				Burst:     10,                   // Allow burst of 10 requests
				ExpiresIn: 3 * time.Minute,      // Clean up expired entries
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
	authRateLimiterConfig := middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      2,                    // 2 requests per second
				Burst:     5,                    // Allow burst of 5 requests
				ExpiresIn: 5 * time.Minute,      // Clean up expired entries
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
	e.Use(middleware.RateLimiterWithConfig(generalRateLimiterConfig))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Public routes with stricter rate limiting
	authGroup := e.Group("")
	authGroup.Use(middleware.RateLimiterWithConfig(authRateLimiterConfig))
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	// Protected routes
	protectedGroup := e.Group("")
	protectedGroup.Use(mw.AuthMiddleware(cfg.JWTSecret, log))
	protectedGroup.GET("/users/me", authHandler.GetCurrentUser)

	// Admin routes (example for future use)
	adminGroup := e.Group("/admin")
	adminGroup.Use(mw.AuthMiddleware(cfg.JWTSecret, log))
	adminGroup.Use(mw.SuperAdminMiddleware(log))
	// Add admin routes here in the future

	log.Info("Starting server on port " + cfg.Port)
	log.Info("Available endpoints:")
	log.Info("  GET  /health")
	log.Info("  POST /register")
	log.Info("  POST /login")
	log.Info("  GET  /users/me (protected)")

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