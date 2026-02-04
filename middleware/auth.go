package middleware

import (
	"github.com/labstack/echo/v4"
	"oauth2-provider/utils"
	"strings"
)

// JWTAuth middleware for API endpoints requiring Bearer token
func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				claims, err := utils.ValidateJWT(token)
				if err == nil {
					c.Set("user_id", claims.Subject)
					return next(c)
				}
			}
		}

		// Also check cookie for convenience in some API calls, though standard is Bearer
		cookie, err := c.Cookie("auth_token")
		if err == nil {
			claims, err := utils.ValidateJWT(cookie.Value)
			if err == nil {
				c.Set("user_id", claims.Subject)
				return next(c)
			}
		}

		return echo.ErrUnauthorized
	}
}

// UserSession middleware checks for auth_token cookie and sets user_id in context if valid.
// It does NOT return an error if authentication fails, allowing the handler to decide.
func UserSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("auth_token")
		if err == nil {
			claims, err := utils.ValidateJWT(cookie.Value)
			if err == nil {
				c.Set("user_id", claims.Subject)
			}
		}
		return next(c)
	}
}
