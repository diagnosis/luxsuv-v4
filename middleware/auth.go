package middleware

import (
	"context"
	"luxsuv-v4/models"
	"luxsuv-v4/services"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	authService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware validates JWT token and adds user info to context
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.sendErrorResponse(w, "authorization header required", http.StatusUnauthorized)
			return
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.sendErrorResponse(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			m.sendErrorResponse(w, "token required", http.StatusUnauthorized)
			return
		}

		// Validate JWT token
		claims, err := m.authService.ValidateJWT(tokenString)
		if err != nil {
			m.sendErrorResponse(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract user information from claims
		userID, ok := (*claims)["user_id"]
		if !ok {
			m.sendErrorResponse(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		username, _ := (*claims)["username"].(string)
		email, _ := (*claims)["email"].(string)
		role, _ := (*claims)["role"].(string)
		superAdmin, _ := (*claims)["super_admin"].(bool)

		// Add user information to request context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "username", username)
		ctx = context.WithValue(ctx, "email", email)
		ctx = context.WithValue(ctx, "role", role)
		ctx = context.WithValue(ctx, "super_admin", superAdmin)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// RequireRole middleware checks if user has required role
func (m *AuthMiddleware) RequireRole(role string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return m.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Context().Value("role")
			superAdmin := r.Context().Value("super_admin")

			// Super admins can access everything
			if superAdmin == true {
				next.ServeHTTP(w, r)
				return
			}

			// Check if user has required role
			if userRole != role {
				m.sendErrorResponse(w, "insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin middleware checks if user is admin or super admin
func (m *AuthMiddleware) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return m.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Context().Value("role")
		superAdmin := r.Context().Value("super_admin")

		if superAdmin == true || userRole == "admin" {
			next.ServeHTTP(w, r)
			return
		}

		m.sendErrorResponse(w, "admin access required", http.StatusForbidden)
	})
}

// sendErrorResponse sends a standardized error response
func (m *AuthMiddleware) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := models.ErrorResponse{Error: message}
	// Simple JSON encoding without importing encoding/json in middleware
	w.Write([]byte(`{"error":"` + message + `"}`))
}