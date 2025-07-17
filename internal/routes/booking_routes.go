package routes

import (
	"github.com/diagnosis/luxsuv-v4/internal/handlers"
	"github.com/diagnosis/luxsuv-v4/internal/middleware"
	"github.com/labstack/echo/v4"
)

// SetupBookingRoutes configures all booking-related routes
func SetupBookingRoutes(e *echo.Echo, bookRideHandler *handlers.BookRideHandler, authMiddleware *middleware.AuthMiddleware) {
	// Public booking routes (no authentication required)
	publicBookingGroup := e.Group("/bookings")
	
	// Create booking (supports both authenticated and guest users)
	e.POST("/book-ride", bookRideHandler.Create, authMiddleware.OptionalAuth())
	
	// Guest booking management
	publicBookingGroup.GET("/email/:email", bookRideHandler.GetByEmail)
	publicBookingGroup.POST("/:id/update-link", bookRideHandler.GenerateUpdateLink)
	
	// Public update/cancel with secure token (for guest users)
	publicBookingGroup.PUT("/:id/update", bookRideHandler.Update)
	publicBookingGroup.DELETE("/:id/cancel", bookRideHandler.Cancel)

	// Protected booking routes (require authentication)
	protectedBookingGroup := e.Group("/bookings")
	protectedBookingGroup.Use(authMiddleware.RequireAuth())
	
	// Authenticated user booking management
	protectedBookingGroup.GET("/my", bookRideHandler.GetByUserID)
	
	// Note: These routes are also handled by the public routes above with token validation
	// but we keep them here for authenticated users who don't need tokens
	protectedBookingGroup.PUT("/:id", bookRideHandler.Update)
	protectedBookingGroup.DELETE("/:id/cancel", bookRideHandler.Cancel)

	// Driver-specific routes
	driverGroup := e.Group("/driver")
	driverGroup.Use(authMiddleware.RequireAuth())
	driverGroup.Use(authMiddleware.RequireDriver())
	
	driverGroup.PUT("/bookings/:id/accept", bookRideHandler.Accept)
}