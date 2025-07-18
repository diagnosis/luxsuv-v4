package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/diagnosis/luxsuv-v4/internal/auth"
	"github.com/diagnosis/luxsuv-v4/internal/email"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/diagnosis/luxsuv-v4/internal/models"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService  *auth.Service
	emailService *email.Service
	logger       *logger.Logger
}

func NewAuthHandler(authService *auth.Service, emailService *email.Service, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		emailService: emailService,
		logger:       logger,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c echo.Context) error {
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Warn(fmt.Sprintf("Invalid request body: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	user, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Registration failed: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	h.logger.Info(fmt.Sprintf("User registered successfully: %s", user.Email))

	// Send welcome email if email service is configured
	if h.emailService != nil {
		if err := h.emailService.SendWelcomeEmail(user.Email, user.Username); err != nil {
			h.logger.Warn(fmt.Sprintf("Failed to send welcome email to %s: %s", user.Email, err.Error()))
			// Don't fail registration if email fails
		}
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "user registered successfully",
		"user":    user,
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Warn(fmt.Sprintf("Invalid request body: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	response, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Login failed: %s", err.Error()))
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	h.logger.Info(fmt.Sprintf("User logged in successfully: %s", response.User.Email))
	return c.JSON(http.StatusOK, response)
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(c echo.Context) error {
	userIDClaim := c.Get("user_id")
	if userIDClaim == nil {
		h.logger.Warn("Missing user_id in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid token",
		})
	}

	var userID int64
	switch v := userIDClaim.(type) {
	case float64:
		userID = int64(v)
	case int64:
		userID = v
	case int:
		userID = int64(v)
	default:
		h.logger.Warn(fmt.Sprintf("Invalid user_id type: %T", userIDClaim))
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid token claims",
		})
	}

	user, err := h.authService.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Failed to get current user: %s", err.Error()))
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser handles user deletion (admin only)
func (h *AuthHandler) DeleteUser(c echo.Context) error {
	// Get admin user ID from context
	adminIDClaim := c.Get("user_id")
	if adminIDClaim == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid token",
		})
	}

	var adminID int64
	switch v := adminIDClaim.(type) {
	case float64:
		adminID = int64(v)
	case int64:
		adminID = v
	case int:
		adminID = int64(v)
	default:
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid token claims",
		})
	}

	// Get user ID from URL parameter
	userIDParam := c.Param("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user ID",
		})
	}

	if err := h.authService.DeleteUser(c.Request().Context(), userID, adminID); err != nil {
		h.logger.Warn(fmt.Sprintf("Failed to delete user: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	h.logger.Info(fmt.Sprintf("User %d deleted by admin %d", userID, adminID))
	return c.JSON(http.StatusOK, map[string]string{
		"message": "user deleted successfully",
	})
}
