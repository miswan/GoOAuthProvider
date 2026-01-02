package middleware

import (
	"github.com/labstack/echo/v4"
	"oauth2-provider/utils"
)

// SessionAuthMiddleware checks for a valid session token in the cookie
func SessionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			// No session cookie, proceed without setting user_id
			// The handler (Authorize) will check if user_id is set and redirect if not
			return next(c)
		}

		token := cookie.Value
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			// Invalid token, treat as unauthenticated
			return next(c)
		}

		c.Set("user_id", claims.Subject)
		return next(c)
	}
}
