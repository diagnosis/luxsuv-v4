package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/diagnosis/luxsuv-v4/internal/auth"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	authService *auth.Service
	logger      *logger.Logger
}

func NewAuthMiddleware(authService *auth.Service, logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// Helper function to convert interface{} to int64
func convertToInt64(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	case float32:
		return int64(v), true
	default:
		return 0, false
	}
}

// RequireAuth middleware validates JWT tokens
func (m *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				m.logger.Warn("Missing Authorization header")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing authorization header",
				})
			}

			// Check Bearer prefix
			if !strings.HasPrefix(authHeader, "Bearer ") {
				m.logger.Warn("Invalid Authorization header format")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization header format",
				})
			}

			// Extract token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				m.logger.Warn("Empty token")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing token",
				})
			}

			// Validate token
			claims, err := m.authService.ValidateJWT(tokenString)
			if err != nil {
				m.logger.Warn(fmt.Sprintf("Token validation failed: %s", err.Error()))
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid token",
				})
			}

			// Extract user information from claims
			userIDRaw, ok := claims["user_id"]
			if !ok {
				m.logger.Warn("Missing user_id in token claims")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid token claims",
				})
			}

			// Convert user_id to int64
			userID, ok := convertToInt64(userIDRaw)
			if !ok {
				m.logger.Warn(fmt.Sprintf("Invalid user_id type in token claims: %T", userIDRaw))
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid token claims",
				})
			}

			role, ok := claims["role"]
			if !ok {
				m.logger.Warn("Missing role in token claims")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid token claims",
				})
			}

			isAdmin, ok := claims["is_admin"]
			if !ok {
				isAdmin = false // Default to false if not present
			}

			// Set context values
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("is_admin", isAdmin)
			c.Set("username", claims["username"])
			c.Set("email", claims["email"])

			m.logger.Info(fmt.Sprintf("Authenticated user: %v (type: %T, role: %v)", userID, userID, role))
			return next(c)
		}
	}
}

// RequireAdmin middleware ensures only admins can access the endpoint
func (m *AuthMiddleware) RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isAdmin := c.Get("is_admin")
			if isAdmin == nil {
				m.logger.Warn("Missing is_admin in context")
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "access denied",
				})
			}

			adminFlag, ok := isAdmin.(bool)
			if !ok || !adminFlag {
				m.logger.Warn("Access denied: user is not an admin")
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "admin access required",
				})
			}

			return next(c)
		}
	}
}

// RequireDriver middleware ensures only drivers can access the endpoint
func (m *AuthMiddleware) RequireDriver() echo.MiddlewareFunc {
	return m.RequireRole("driver")
}

// RequireRole middleware ensures user has specific role
func (m *AuthMiddleware) RequireRole(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get("role")
			if role == nil {
				m.logger.Warn("Missing role in context")
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "access denied",
				})
			}

			userRole, ok := role.(string)
			if !ok {
				m.logger.Warn("Invalid role type in context")
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "access denied",
				})
			}

			// Check if user has required role or is admin
			isAdmin := c.Get("is_admin")
			adminFlag, _ := isAdmin.(bool)

			if userRole != requiredRole && !adminFlag {
				m.logger.Warn(fmt.Sprintf("Access denied: required role %s, user role %s", requiredRole, userRole))
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "insufficient permissions",
				})
			}

			return next(c)
		}
	}
}

// OptionalAuth middleware sets user context if token present, but doesn't require it
func (m *AuthMiddleware) OptionalAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return next(c) // No token, proceed as guest
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				return next(c)
			}

			claims, err := m.authService.ValidateJWT(tokenString)
			if err != nil {
				return next(c) // Invalid token, proceed as guest
			}

			userIDRaw, ok := claims["user_id"]
			if !ok {
				return next(c)
			}

			// Convert user_id to int64
			userID, ok := convertToInt64(userIDRaw)
			if !ok {
				return next(c)
			}

			role, ok := claims["role"]
			if !ok {
				return next(c)
			}

			isAdmin, ok := claims["is_admin"]
			if !ok {
				isAdmin = false
			}

			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("is_admin", isAdmin)
			c.Set("username", claims["username"])
			c.Set("email", claims["email"])

			m.logger.Info(fmt.Sprintf("Authenticated user (optional): %v (type: %T, role: %v)", userID, userID, role))
			return next(c)
		}
	}
}