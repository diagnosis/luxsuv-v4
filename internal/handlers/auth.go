package handlers

import (
	"fmt"
	"github.com/diagnosis/luxsuv-v4/internal/data"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/diagnosis/luxsuv-v4/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"net/http"
)

type AuthHandler struct {
	authService *services.AuthService
	log         *logger.Logger
}

func NewAuthHandler(authService *services.AuthService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, log: log}
}

func (h *AuthHandler) Register(c echo.Context) error {
	user := new(data.User)
	if err := c.Bind(user); err != nil {
		h.log.Err("Invalid request body for register: " + err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err := h.authService.Register(user)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			h.log.Warn("Registration failed: email already exists: " + user.Email)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "email already exists"})
		}
		h.log.Warn("Registration failed: " + err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		h.log.Err("Invalid request body for login: " + err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) GetCurrentUser(c echo.Context) error {
	userID, ok := c.Get("user_id").(float64)
	if !ok {
		h.log.Err("Invalid user ID type in JWT context: " + fmt.Sprintf("%T", c.Get("user_id")))
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user ID type"})
	}

	user, err := h.authService.GetUserByID(int64(userID))
	if err != nil {
		h.log.Err("Failed to get current user: " + err.Error())
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, user)
}
