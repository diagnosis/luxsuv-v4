package routes

import (
	"github.com/diagnosis/luxsuv-v4/internal/handlers"
	"github.com/diagnosis/luxsuv-v4/internal/middleware"
	"github.com/labstack/echo/v4"
)

// SetupAdminRoutes configures all admin-related routes
func SetupAdminRoutes(e *echo.Echo, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	// Admin routes - require authentication and admin role
	adminGroup := e.Group("/admin")
	adminGroup.Use(authMiddleware.RequireAuth())
	adminGroup.Use(authMiddleware.RequireAdmin())

	// User management endpoints
	adminGroup.GET("/users", userHandler.ListUsers)
	adminGroup.GET("/users/by-email", userHandler.GetUserByEmail)
	adminGroup.GET("/users/:id", userHandler.GetUserByID)
	adminGroup.PUT("/users/:id/role", userHandler.UpdateUserRole)
	adminGroup.DELETE("/users/:id", authHandler.DeleteUser)
}