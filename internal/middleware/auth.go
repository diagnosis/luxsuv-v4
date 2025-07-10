package mw

import (
	"fmt"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func AuthMiddleware(jwtSecret string, log *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				log.Warn("Auth middleware: missing or invalid token")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing or invalid token"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				log.Warn("Auth middleware: invalid token: " + err.Error())
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Warn("Auth middleware: invalid token claims")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
			}

			c.Set("user_id", claims["id"])
			c.Set("role", claims["role"])
			c.Set("super_admin", claims["super_admin"])
			log.Info("Auth middleware: successfully validated token for user ID " + fmt.Sprintf("%v", claims["id"]))
			return next(c)
		}
	}
}
