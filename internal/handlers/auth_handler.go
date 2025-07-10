package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/diagnosis/luxsuv-v4/internal/auth"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/diagnosis/luxsuv-v4/internal/models"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *auth.Service
	logger      *logger.Logger
}

func NewAuthHandler(authService *auth.Service, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
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

	h.logger.Info(fmt.Sprintf("Registration request received: username=%s, email=%s, role=%s", req.Username, req.Email, req.Role))

	user, err := h.authService.Register(&req)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Registration failed: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	h.logger.Info(fmt.Sprintf("User registered successfully: %s", user.Email))
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

	response, err := h.authService.Login(&req)
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

	user, err := h.authService.GetUserByID(userID)
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

	if err := h.authService.DeleteUser(userID, adminID); err != nil {
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