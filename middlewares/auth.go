package middlewares

import (
	"net/http"
	"time"

	"notification-server/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenHeader := c.Request().Header.Get("Authorization")
		if tokenHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
		}

		tokenString := tokenHeader[len("Bearer "):]

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetEnv("JWT_SECRET")), nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			if claims.ExpiresAt.Time.Before(time.Now()) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token expired"})
			}
			c.Set("userID", claims.UserID)
			return next(c)
		}

		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}
}
