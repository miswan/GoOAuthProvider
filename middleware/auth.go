package middleware

import (
	"github.com/labstack/echo/v4"
	"oauth2-provider/utils"
	"strings"
)

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.ErrUnauthorized
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.ErrUnauthorized
		}

		token := parts[1]
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			return echo.ErrUnauthorized
		}

		c.Set("user_id", claims.Subject)
		return next(c)
	}
}

// SessionAuthMiddleware checks for a session token in cookies
func SessionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			// No session cookie, user is not authenticated
			// We don't redirect here because this middleware might be used in a chain where we just want to know IF the user is authenticated.
			// Handlers should check c.Get("user_id")
			return next(c)
		}

		claims, err := utils.ValidateJWT(cookie.Value)
		if err != nil {
			// Invalid token
			return next(c)
		}

		c.Set("user_id", claims.Subject)
		return next(c)
	}
}
