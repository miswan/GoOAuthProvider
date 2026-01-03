package middleware

import (
	"github.com/labstack/echo/v4"
	"oauth2-provider/utils"
)

func SessionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			// No session cookie, proceed without user_id (unauthenticated)
			return next(c)
		}

		claims, err := utils.ValidateJWT(cookie.Value)
		if err != nil {
			// Invalid token, ignore it
			return next(c)
		}

		// Set user_id in context
		c.Set("user_id", claims.Subject)
		return next(c)
	}
}

func RequireLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("user_id")
		if userID == nil {
			// Redirect to login page
			// We can pass the current URL as redirect_to
			return c.Redirect(302, "/login?redirect_to="+c.Request().URL.String())
		}
		return next(c)
	}
}
