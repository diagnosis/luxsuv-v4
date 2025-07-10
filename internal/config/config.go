package config

import (
	"errors"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

func LoadConfig(log *logger.Logger) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Err("Failed to load .env file: " + err.Error())
		return nil, err
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Port:        os.Getenv("PORT"),
	}

	if cfg.DatabaseURL == "" || cfg.JWTSecret == "" || cfg.Port == "" {
		log.Err("Missing required environment variables")
		return nil, errors.New("missing required environment variables")
	}

	log.Info("Loaded DATABASE_URL: " + cfg.DatabaseURL)
	return cfg, nil
}
