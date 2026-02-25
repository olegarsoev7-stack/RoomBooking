package auth

import (
	"net/http"
	"strings"

	"search-job/pkg/jwt"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware защищает роуты
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
		}

		// Заголовок должен быть в формате "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization format"})
		}

		tokenString := parts[1]
		userID, err := jwt.ParseToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
		}

		// Сохраняем UserID в контекст, чтобы потом использовать в обработчиках (например, при создании брони)
		c.Set("userID", userID)

		return next(c)
	}
}
