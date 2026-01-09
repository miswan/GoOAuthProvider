package middleware

import (
	"github.com/labstack/echo/v4"
	"oauth2-provider/utils"
)

func SessionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			return next(c)
		}

		// Verify session token (JWT)
		claims, err := utils.ValidateJWT(cookie.Value)
		if err != nil {
			// Invalid session, proceed without user_id
			return next(c)
		}

		c.Set("user_id", claims.Subject)
		return next(c)
	}
}
