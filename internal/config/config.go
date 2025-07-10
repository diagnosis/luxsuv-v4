package config

import (
	"errors"
	"os"
	"strconv"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	Port           string
	Environment    string
	LogLevel       string
	MaxConnections int
}

func LoadConfig(log *logger.Logger) (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist (for production)
	if err := godotenv.Load(); err != nil {
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

	return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}