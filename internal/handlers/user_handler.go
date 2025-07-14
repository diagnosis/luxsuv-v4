package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/diagnosis/luxsuv-v4/internal/auth"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/diagnosis/luxsuv-v4/internal/models"
	"github.com/diagnosis/luxsuv-v4/internal/repository"
	"github.com/diagnosis/luxsuv-v4/internal/validation"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	authService *auth.Service
	userRepo    repository.UserRepository
	logger      *logger.Logger
}

func NewUserHandler(authService *auth.Service, userRepo repository.UserRepository, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		authService: authService,
		userRepo:    userRepo,
		logger:      logger,
	}
}

// ListUsers handles listing all users with pagination (admin only)
func (h *UserHandler) ListUsers(c echo.Context) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 10 // Default limit
	}

	offset := (page - 1) * limit

	// Get users from database
	users, err := h.userRepo.ListUsers(c.Request().Context(), limit, offset)
	if err != nil {
		h.logger.Err(fmt.Sprintf("Failed to list users: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to retrieve users",
		})
	}

	// Get total count for pagination
	totalCount, err := h.userRepo.CountUsers(c.Request().Context())
	if err != nil {
		h.logger.Err(fmt.Sprintf("Failed to count users: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to retrieve user count",
		})
	}

	totalPages := (totalCount + int64(limit) - 1) / int64(limit)

	response := map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"current_page": page,
			"total_pages":  totalPages,
			"total_count":  totalCount,
			"limit":        limit,
		},
	}

	h.logger.Info(fmt.Sprintf("Listed %d users (page %d)", len(users), page))
	return c.JSON(http.StatusOK, response)
}

// GetUserByEmail handles retrieving a user by email (admin only)
func (h *UserHandler) GetUserByEmail(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "email parameter is required",
		})
	}

	// Validate email format
	email = strings.TrimSpace(strings.ToLower(email))
	if err := validation.ValidateEmail(email); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	user, err := h.userRepo.GetByEmail(c.Request().Context(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "user not found",
			})
		}
		h.logger.Err(fmt.Sprintf("Failed to get user by email %s: %s", email, err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to retrieve user",
		})
	}

	// Remove password from response
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

// GetUserByID handles retrieving a specific user by ID (admin only)
func (h *UserHandler) GetUserByID(c echo.Context) error {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user ID",
		})
	}

	user, err := h.authService.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Failed to get user by ID %d: %s", userID, err.Error()))
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUserRole handles updating a user's role (admin only)
func (h *UserHandler) UpdateUserRole(c echo.Context) error {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user ID",
		})
	}

	var req struct {
		Role string `json:"role" validate:"required,oneof=rider driver admin"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// Validate role
	if !models.IsValidRole(req.Role) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid role; must be rider, driver, or admin",
		})
	}

	// Get the user to update
	user, err := h.userRepo.GetByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
	}

	// Update user role
	user.Role = req.Role
	user.IsAdmin = req.Role == models.RoleAdmin

	if err := h.userRepo.UpdateUserRole(c.Request().Context(), userID, req.Role, user.IsAdmin); err != nil {
		h.logger.Err(fmt.Sprintf("Failed to update user role: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update user role",
		})
	}

	h.logger.Info(fmt.Sprintf("User role updated: ID %d, new role: %s", userID, req.Role))

	// Remove password from response
	user.Password = ""
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "user role updated successfully",
		"user":    user,
	})
}
