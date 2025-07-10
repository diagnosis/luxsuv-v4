package main

import (
	"fmt"
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
	log.Info("Loaded DATABASE_URL: " + cfg.DatabaseURL)

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

	// Connect to database
	db, err = sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Err("Failed to connect to database: " + err.Error())
		return
	}
	defer db.Close()
	log.Info("Successfully connected to database")

	// Initialize dependencies
	repo := data.NewRepository(db)
	authService := services.NewAuthService(repo, cfg.JWTSecret, log)
	authHandler := handlers.NewAuthHandler(authService, log)

	// Set up Echo server
	e := echo.New()
	// Configure rate limiter with stricter limits
	rateLimiterConfig := middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(middleware.RateLimiterMemoryStoreRate(10)), // 10 requests per second
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
	}
	loginRateLimiterConfig := middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(middleware.RateLimiterMemoryStoreRate(5)), // 5 login attempts per second
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
	}
	e.Use(middleware.RateLimiter(rateLimiterConfig))
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login, middleware.RateLimiter(loginRateLimiterConfig))
	e.GET("/users/me", authHandler.GetCurrentUser, mw.AuthMiddleware(cfg.JWTSecret, log))

	log.Info("Starting server on port " + cfg.Port)
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
