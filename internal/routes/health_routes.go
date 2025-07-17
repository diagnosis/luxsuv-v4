package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// SetupHealthRoutes configures health check and system status routes
func SetupHealthRoutes(e *echo.Echo) {
	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API info endpoint (optional)
	e.GET("/api/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"service": "LuxSUV Backend API",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": map[string]interface{}{
				"auth": []string{
					"POST /register",
					"POST /login",
					"POST /auth/forgot-password",
					"POST /auth/reset-password",
					"GET /users/me",
					"PUT /users/me/password",
				},
				"bookings": []string{
					"POST /book-ride",
					"GET /bookings/email/:email",
					"POST /bookings/:id/update-link",
					"GET /bookings/my",
					"PUT /bookings/:id",
					"DELETE /bookings/:id/cancel",
					"PUT /driver/bookings/:id/accept",
				},
				"admin": []string{
					"GET /admin/users",
					"GET /admin/users/by-email",
					"GET /admin/users/:id",
					"PUT /admin/users/:id/role",
					"DELETE /admin/users/:id",
				},
			},
		})
	})
}