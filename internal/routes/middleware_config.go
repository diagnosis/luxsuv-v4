package routes

import (
	"time"

	"github.com/diagnosis/luxsuv-v4/internal/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// MiddlewareConfig holds all middleware configurations
type MiddlewareConfig struct {
	GeneralRateLimiter echomiddleware.RateLimiterConfig
	AuthRateLimiter    echomiddleware.RateLimiterConfig
}

// SetupGlobalMiddleware configures global middleware for the Echo instance
func SetupGlobalMiddleware(e *echo.Echo, environment string) MiddlewareConfig {
	// Basic middleware
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.Logger())

	// CORS configuration based on environment
	if environment == "development" {
		e.Use(echomiddleware.CORSWithConfig(middleware.DevelopmentCORSConfig()))
	} else {
		e.Use(echomiddleware.CORSWithConfig(middleware.CORSConfig()))
	}

	e.Use(echomiddleware.Secure())

	// Rate limiter configurations
	generalRateLimiterConfig := echomiddleware.RateLimiterConfig{
		Store: echomiddleware.NewRateLimiterMemoryStoreWithConfig(
			echomiddleware.RateLimiterMemoryStoreConfig{
				Rate:      5,               // 5 requests per second
				Burst:     10,              // Allow burst of 10 requests
				ExpiresIn: 3 * time.Minute, // Clean up expired entries
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(429, map[string]string{"error": "too many requests"})
		},
	}

	authRateLimiterConfig := echomiddleware.RateLimiterConfig{
		Store: echomiddleware.NewRateLimiterMemoryStoreWithConfig(
			echomiddleware.RateLimiterMemoryStoreConfig{
				Rate:      2,               // 2 requests per second
				Burst:     5,               // Allow burst of 5 requests
				ExpiresIn: 5 * time.Minute, // Clean up expired entries
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(429, map[string]string{"error": "too many requests"})
		},
	}

	// Apply general rate limiting globally
	e.Use(echomiddleware.RateLimiterWithConfig(generalRateLimiterConfig))

	return MiddlewareConfig{
		GeneralRateLimiter: generalRateLimiterConfig,
		AuthRateLimiter:    authRateLimiterConfig,
	}
}