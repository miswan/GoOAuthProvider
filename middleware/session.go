package middleware

import (
	"github.com/labstack/echo/v4"
	"oauth2-provider/utils"
	"strings"
)

func SessionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Try to get token from Authorization header first
		authHeader := c.Request().Header.Get("Authorization")
		var tokenString string

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		} else {
			// Try to get from cookie
			cookie, err := c.Cookie("session_token")
			if err == nil {
				tokenString = cookie.Value
			}
		}

		if tokenString != "" {
			claims, err := utils.ValidateJWT(tokenString)
			if err == nil {
				c.Set("user_id", claims.Subject)
			}
		}

		return next(c)
	}
}
