package config

import (
	"errors"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	Port           string
	Environment    string
	LogLevel       string
	MaxConnections int

	// Email configuration (MailerSend)
	MailerSendAPIKey    string
	MailerSendFromEmail string
	MailerSendFromName  string
}

func LoadConfig(log *logger.Logger) (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist (for production)
	if err := godotenv.Load("../../.env"); err != nil {
		log.Warn("No .env file found, using environment variables: " + err.Error())
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Port:        getEnvWithDefault("PORT", "8080"),
		Environment: getEnvWithDefault("ENVIRONMENT", "development"),
		LogLevel:    getEnvWithDefault("LOG_LEVEL", "info"),
	}

	// Parse max connections
	maxConnStr := getEnvWithDefault("MAX_DB_CONNECTIONS", "25")
	maxConn, err := strconv.Atoi(maxConnStr)
	if err != nil {
		log.Warn("Invalid MAX_DB_CONNECTIONS value, using default: " + err.Error())
		maxConn = 25
	}
	cfg.MaxConnections = maxConn

	// Email configuration (MailerSend)
	cfg.MailerSendAPIKey = getEnvWithDefault("MAILERSEND_API_KEY", "")
	cfg.MailerSendFromEmail = getEnvWithDefault("MAILERSEND_FROM_EMAIL", "")
	cfg.MailerSendFromName = getEnvWithDefault("MAILERSEND_FROM_NAME", "LuxSUV Support")

	// Validate required fields
	if cfg.DatabaseURL == "" {
		log.Err("DATABASE_URL environment variable is required")
		return nil, errors.New("DATABASE_URL is required")
	}

	if cfg.JWTSecret == "" {
		log.Err("JWT_SECRET environment variable is required")
		return nil, errors.New("JWT_SECRET is required")
	}

	// Validate JWT secret strength
	if len(cfg.JWTSecret) < 32 {
		log.Err("JWT_SECRET must be at least 32 characters long")
		return nil, errors.New("JWT_SECRET must be at least 32 characters long")
	}

	log.Info("Configuration loaded successfully")
	log.Info("Environment: " + cfg.Environment)
	log.Info("Port: " + cfg.Port)
	log.Info("Log Level: " + cfg.LogLevel)
	log.Info("Max DB Connections: " + strconv.Itoa(cfg.MaxConnections))

	// Log email configuration (without sensitive API key)
	if cfg.MailerSendAPIKey != "" {
		log.Info("MailerSend From Email: " + cfg.MailerSendFromEmail)
		log.Info("MailerSend From Name: " + cfg.MailerSendFromName)
		log.Info("Email service enabled")
	} else {
		log.Warn("Email service not configured - MAILERSEND_API_KEY not set")
	}

	return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
