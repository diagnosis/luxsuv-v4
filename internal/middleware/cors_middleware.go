package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CORSConfig returns a comprehensive CORS configuration
func CORSConfig() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",    // React dev server
			"http://localhost:3001",    // Alternative React port
			"http://localhost:5173",    // Vite dev server
			"http://localhost:8080",    // Local development
			"http://127.0.0.1:3000",   // Alternative localhost
			"http://127.0.0.1:5173",   // Alternative localhost for Vite
			// Add your production domains here
			// "https://yourdomain.com",
			// "https://www.yourdomain.com",
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.PATCH,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-CSRF-Token",
			"Cache-Control",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Total-Count",
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}
}

// DevelopmentCORSConfig returns a permissive CORS config for development
func DevelopmentCORSConfig() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false, // Cannot be true when AllowOrigins is "*"
		MaxAge:           86400,
	}
}