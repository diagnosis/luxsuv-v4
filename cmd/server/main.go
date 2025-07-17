package main

import (
	"fmt"
	"time"

	"github.com/diagnosis/luxsuv-v4/internal/auth"
	"github.com/diagnosis/luxsuv-v4/internal/config"
	"github.com/diagnosis/luxsuv-v4/internal/email"
	"github.com/diagnosis/luxsuv-v4/internal/handlers"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/diagnosis/luxsuv-v4/internal/middleware"
	"github.com/diagnosis/luxsuv-v4/internal/repository/postgres"
	"github.com/diagnosis/luxsuv-v4/internal/routes"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
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

	// Initialize database
	db, err := initializeDatabase(cfg, log)
	if err != nil {
		log.Err("Failed to initialize database: " + err.Error())
		return
	}
	defer db.Close()

	// Initialize services
	services, err := initializeServices(db, cfg, log)
	if err != nil {
		log.Err("Failed to initialize services: " + err.Error())
		return
	}

	// Initialize handlers
	handlers := initializeHandlers(services, log)

	// Set up Echo server
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Setup global middleware
	middlewareConfig := routes.SetupGlobalMiddleware(e, cfg.Environment)
	log.Info(fmt.Sprintf("Using %s CORS configuration", cfg.Environment))

	// Setup all routes
	setupAllRoutes(e, handlers, services.AuthMiddleware, middlewareConfig)

	// Log available endpoints
	logAvailableEndpoints(log)

	// Start server
	log.Info("Starting server on port " + cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Err("Failed to start server: " + err.Error())
	}
}

// Services holds all initialized services
type Services struct {
	AuthService    *auth.Service
	EmailService   *email.Service
	AuthMiddleware *middleware.AuthMiddleware
}

// Handlers holds all initialized handlers
type Handlers struct {
	AuthHandler     *handlers.AuthHandler
	UserHandler     *handlers.UserHandler
	PasswordHandler *handlers.PasswordHandler
	BookRideHandler *handlers.BookRideHandler
}

// initializeDatabase sets up database connection and runs migrations
func initializeDatabase(cfg *config.Config, log *logger.Logger) (*sqlx.DB, error) {
	// Run migrations
	migrationDB, err := sqlx.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database for migrations: %w", err)
	}
	defer migrationDB.Close()

	goose.SetLogger(&GooseLogger{log: log})
	if err := goose.Up(migrationDB.DB, "../../migrations"); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}
	log.Info("Database migrations applied successfully")

	// Connect to database with connection pool settings
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Info("Successfully connected to database")

	return db, nil
}

// initializeServices creates and configures all services
func initializeServices(db *sqlx.DB, cfg *config.Config, log *logger.Logger) (*Services, error) {
	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	
	// Initialize auth service
	authService := auth.NewService(userRepo, cfg.JWTSecret, log)
	
	// Initialize email service
	var emailService *email.Service
	if cfg.MailerSendAPIKey != "" && cfg.MailerSendFromEmail != "" {
		emailConfig := email.Config{
			APIKey:    cfg.MailerSendAPIKey,
			FromEmail: cfg.MailerSendFromEmail,
			FromName:  cfg.MailerSendFromName,
		}
		emailService = email.NewService(emailConfig, log)
		log.Info("Email service initialized")
		log.Info("MailerSend From Email: " + cfg.MailerSendFromEmail)
		log.Info("MailerSend From Name: " + cfg.MailerSendFromName)
	} else {
		log.Warn("Email service disabled - MailerSend configuration incomplete")
		log.Warn("Please configure MAILERSEND_API_KEY and MAILERSEND_FROM_EMAIL in .env file")
	}

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService, log)

	return &Services{
		AuthService:    authService,
		EmailService:   emailService,
		AuthMiddleware: authMiddleware,
	}, nil
}

// initializeHandlers creates all handlers
func initializeHandlers(services *Services, log *logger.Logger) *Handlers {
	// Get database connection for repositories
	// TODO: Refactor to pass repositories through services for better architecture
	cfg, err := config.LoadConfig(log)
	if err != nil {
		log.Err("Failed to load config for handlers: " + err.Error())
		panic("Failed to initialize handlers")
	}
	
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Err("Failed to connect to database for handlers: " + err.Error())
		panic("Failed to initialize handlers")
	}
	
	userRepo := postgres.NewUserRepository(db)
	bookRideRepo := postgres.NewBookRideRepository(db)

	return &Handlers{
		AuthHandler:     handlers.NewAuthHandler(services.AuthService, services.EmailService, log),
		UserHandler:     handlers.NewUserHandler(services.AuthService, userRepo, log),
		PasswordHandler: handlers.NewPasswordHandler(services.AuthService, userRepo, services.EmailService, log),
		BookRideHandler: handlers.NewBookRideHandler(bookRideRepo, log, services.AuthService, services.EmailService),
	}
}

// setupAllRoutes configures all application routes
func setupAllRoutes(e *echo.Echo, handlers *Handlers, authMiddleware *middleware.AuthMiddleware, middlewareConfig routes.MiddlewareConfig) {
	// Health and system routes
	routes.SetupHealthRoutes(e)

	// Authentication routes
	routes.SetupAuthRoutes(e, handlers.AuthHandler, handlers.PasswordHandler, authMiddleware, middlewareConfig.AuthRateLimiter)

	// Admin routes
	routes.SetupAdminRoutes(e, handlers.AuthHandler, handlers.UserHandler, authMiddleware)

	// Booking routes
	routes.SetupBookingRoutes(e, handlers.BookRideHandler, authMiddleware)
}

// logAvailableEndpoints logs all available API endpoints
func logAvailableEndpoints(log *logger.Logger) {
	log.Info("Available endpoints:")
	
	// Health endpoints
	log.Info("  GET  /health")
	log.Info("  GET  /api/info")
	
	// Auth endpoints
	log.Info("  POST /register")
	log.Info("  POST /login")
	log.Info("  POST /auth/forgot-password")
	log.Info("  POST /auth/reset-password")
	log.Info("  GET  /users/me (protected)")
	log.Info("  PUT  /users/me/password (protected)")
	
	// Admin endpoints
	log.Info("  GET  /admin/users (admin only)")
	log.Info("  GET  /admin/users/by-email?email=user@example.com (admin only)")
	log.Info("  GET  /admin/users/:id (admin only)")
	log.Info("  PUT  /admin/users/:id/role (admin only)")
	log.Info("  DELETE /admin/users/:id (admin only)")
	
	// Booking endpoints
	log.Info("  POST /book-ride (public)")
	log.Info("  GET  /bookings/email/:email (public)")
	log.Info("  POST /bookings/:id/update-link (public)")
	log.Info("  GET  /bookings/my (protected)")
	log.Info("  PUT  /bookings/:id (protected/token)")
	log.Info("  DELETE /bookings/:id/cancel (protected/token)")
	log.Info("  PUT  /driver/bookings/:id/accept (driver only)")
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