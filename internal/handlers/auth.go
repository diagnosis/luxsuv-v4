package handlers

import (
	"database/sql"
	"fmt"
	"github.com/diagnosis/luxsuv-v4/internal/data"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/diagnosis/luxsuv-v4/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"net/http"
	"strings"
)

type AuthHandler struct {
	authService *services.AuthService
	log         *logger.Logger
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"omitempty,oneof=admin driver rider"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(authService *services.AuthService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, log: log}
}

func (h *AuthHandler) Register(c echo.Context) error {
	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		h.log.Err("Invalid request body for register: " + err.Error())
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	// Convert request to user model
	user := &data.User{
		Username: strings.TrimSpace(req.Username),
		Email:    strings.TrimSpace(strings.ToLower(req.Email)),
		Password: req.Password,
		Role:     strings.TrimSpace(strings.ToLower(req.Role)),
	}

	err := h.authService.Register(user)
	if err != nil {
		// Check for PostgreSQL unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			h.log.Warn(fmt.Sprintf("PostgreSQL error: Code=%s, Message=%s", pqErr.Code, pqErr.Message))
			if pqErr.Code == "23505" { // unique_violation
				if strings.Contains(pqErr.Message, "email") {
					h.log.Warn("Registration failed: email already exists: " + user.Email)
					return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "email already exists"})
				}
				if strings.Contains(pqErr.Message, "username") {
					h.log.Warn("Registration failed: username already exists: " + user.Username)
					return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "username already exists"})
				}
			}
		}

		// Handle validation errors from service
		h.log.Warn("Registration failed: " + err.Error())
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	h.log.Info("User registered successfully: " + user.Email)
	return c.JSON(http.StatusCreated, SuccessResponse{Message: "user registered successfully"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		h.log.Err("Invalid request body for login: " + err.Error())
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
	}

	// Normalize email
	email := strings.TrimSpace(strings.ToLower(req.Email))
	password := req.Password

	if email == "" || password == "" {
		h.log.Warn("Login failed: missing email or password")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "email and password are required"})
	}

	token, err := h.authService.Login(email, password)
	if err != nil {
		h.log.Warn("Login failed for email " + email + ": " + err.Error())
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
	}

	h.log.Info("User logged in successfully: " + email)
	return c.JSON(http.StatusOK, LoginResponse{Token: token})
}

func (h *AuthHandler) GetCurrentUser(c echo.Context) error {
	// Get user ID from JWT claims set by middleware
	userIDClaim := c.Get("user_id")
	if userIDClaim == nil {
		h.log.Err("User ID not found in JWT claims")
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid token"})
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
		h.log.Err(fmt.Sprintf("Invalid user ID type in JWT context: %T, value: %v", userIDClaim, userIDClaim))
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid token claims"})
	}

	if userID <= 0 {
		h.log.Err(fmt.Sprintf("Invalid user ID value: %d", userID))
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user ID"})
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.log.Warn(fmt.Sprintf("User not found for ID: %d", userID))
			return c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		h.log.Err(fmt.Sprintf("Failed to get current user (ID: %d): %s", userID, err.Error()))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}

	// Remove password from response
	user.Password = ""
	
	h.log.Info(fmt.Sprintf("Retrieved current user: %s (ID: %d)", user.Email, user.ID))
	return c.JSON(http.StatusOK, user)
}