package mw

import (
	"fmt"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func AuthMiddleware(jwtSecret string, log *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				log.Warn("Auth middleware: missing Authorization header")
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "missing authorization header"})
			}

			// Check Bearer prefix
			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Warn("Auth middleware: invalid Authorization header format")
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid authorization header format"})
			}

			// Extract token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				log.Warn("Auth middleware: empty token")
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "missing token"})
			}

			// Parse and validate token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				log.Warn("Auth middleware: token parsing error: " + err.Error())
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid token"})
			}

			if !token.Valid {
				log.Warn("Auth middleware: invalid token")
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid token"})
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Warn("Auth middleware: invalid token claims type")
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid token claims"})
			}

			// Validate required claims
			userID, ok := claims["id"]
			if !ok || userID == nil {
				log.Warn("Auth middleware: missing user ID in token claims")
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid token claims"})
			}

			role, ok := claims["role"]
			if !ok || role == nil {
				log.Warn("Auth middleware: missing role in token claims")
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid token claims"})
			}

			superAdmin, ok := claims["super_admin"]
			if !ok {
				// Default to false if not present for backward compatibility
				superAdmin = false
			}

			// Set context values
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("super_admin", superAdmin)
			
			// Optional: set additional claims
			if username, ok := claims["username"]; ok {
				c.Set("username", username)
			}
			if email, ok := claims["email"]; ok {
				c.Set("email", email)
			}

			log.Info(fmt.Sprintf("Auth middleware: successfully validated token for user ID %v (role: %v)", userID, role))
			return next(c)
		}
	}
}

// RoleMiddleware checks if the user has the required role
func RoleMiddleware(requiredRole string, log *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get("role")
			if role == nil {
				log.Warn("Role middleware: missing role in context")
				return c.JSON(http.StatusForbidden, ErrorResponse{Error: "access denied"})
			}

			userRole, ok := role.(string)
			if !ok {
				log.Warn("Role middleware: invalid role type in context")
				return c.JSON(http.StatusForbidden, ErrorResponse{Error: "access denied"})
			}

			// Check if user has required role or is super admin
			superAdmin := c.Get("super_admin")
			isSuperAdmin, _ := superAdmin.(bool)

			if userRole != requiredRole && !isSuperAdmin {
				log.Warn(fmt.Sprintf("Role middleware: insufficient permissions. Required: %s, User: %s, SuperAdmin: %v", requiredRole, userRole, isSuperAdmin))
				return c.JSON(http.StatusForbidden, ErrorResponse{Error: "insufficient permissions"})
			}

			return next(c)
		}
	}
}

// SuperAdminMiddleware checks if the user is a super admin
func SuperAdminMiddleware(log *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			superAdmin := c.Get("super_admin")
			isSuperAdmin, ok := superAdmin.(bool)
			
			if !ok || !isSuperAdmin {
				log.Warn("SuperAdmin middleware: access denied - user is not a super admin")
				return c.JSON(http.StatusForbidden, ErrorResponse{Error: "super admin access required"})
			}

			return next(c)
		}
	}
}