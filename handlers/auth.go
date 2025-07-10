package handlers

import (
	"encoding/json"
	"luxsuv-v4/models"
	"luxsuv-v4/services"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.RegisterUser(req)
	if err != nil {
		// Handle specific validation errors
		switch err.Error() {
		case "invalid email format":
			h.sendErrorResponse(w, "invalid email format", http.StatusBadRequest)
		case "password must be at least 8 characters long":
			h.sendErrorResponse(w, "password must be at least 8 characters long", http.StatusBadRequest)
		case "email already exists":
			h.sendErrorResponse(w, "email already exists", http.StatusConflict)
		case "username already exists":
			h.sendErrorResponse(w, "username already exists", http.StatusConflict)
		default:
			h.sendErrorResponse(w, "failed to register", http.StatusInternalServerError)
		}
		return
	}

	response := models.AuthResponse{
		Message: "user registered successfully",
		User:    user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, user, err := h.authService.LoginUser(req)
	if err != nil {
		switch err.Error() {
		case "invalid email format":
			h.sendErrorResponse(w, "invalid email format", http.StatusBadRequest)
		case "invalid credentials":
			h.sendErrorResponse(w, "invalid credentials", http.StatusUnauthorized)
		default:
			h.sendErrorResponse(w, "login failed", http.StatusInternalServerError)
		}
		return
	}

	response := models.AuthResponse{
		Message: "login successful",
		Token:   token,
		User:    user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCurrentUser returns the current authenticated user's information
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDValue := r.Context().Value("user_id")
	if userIDValue == nil {
		h.sendErrorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var userID int
	switch v := userIDValue.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	case string:
		var err error
		userID, err = strconv.Atoi(v)
		if err != nil {
			h.sendErrorResponse(w, "invalid user ID", http.StatusBadRequest)
			return
		}
	default:
		h.sendErrorResponse(w, "invalid user ID format", http.StatusBadRequest)
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			h.sendErrorResponse(w, "user not found", http.StatusNotFound)
		} else {
			h.sendErrorResponse(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// sendErrorResponse sends a standardized error response
func (h *AuthHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}