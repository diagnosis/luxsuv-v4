package routes

import (
	"github.com/diagnosis/luxsuv-v4/internal/handlers"
	"github.com/diagnosis/luxsuv-v4/internal/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// SetupAuthRoutes configures all authentication-related routes
func SetupAuthRoutes(e *echo.Echo, authHandler *handlers.AuthHandler, passwordHandler *handlers.PasswordHandler, authMiddleware *middleware.AuthMiddleware, authRateLimiterConfig echomiddleware.RateLimiterConfig) {
	// Public auth routes with stricter rate limiting
	authGroup := e.Group("")
	authGroup.Use(echomiddleware.RateLimiterWithConfig(authRateLimiterConfig))
	
	// Registration and login
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	// Password reset routes (public)
	e.POST("/auth/forgot-password", passwordHandler.ResetPasswordRequest)
	e.POST("/auth/reset-password", passwordHandler.ResetPassword)

	// Protected auth routes
	protectedAuthGroup := e.Group("")
	protectedAuthGroup.Use(authMiddleware.RequireAuth())
	
	// User profile and password management
	protectedAuthGroup.GET("/users/me", authHandler.GetCurrentUser)
	protectedAuthGroup.PUT("/users/me/password", passwordHandler.ChangePassword)
}