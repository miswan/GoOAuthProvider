package middleware

import (
	"github.com/labstack/echo/v4"
	"oauth2-provider/utils"
	"net/http"
)

func SessionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			return next(c) // Continue without setting user_id
		}

		claims, err := utils.ValidateJWT(cookie.Value)
		if err != nil {
			return next(c) // Continue without setting user_id
		}

		c.Set("user_id", claims.Subject)
		return next(c)
	}
}

func RequireLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("user_id")
		if userID == nil {
			// If not authenticated, redirect to login
			// Preserve the current URL as redirect_to
			redirectURL := c.Request().RequestURI
			return c.Redirect(http.StatusFound, "/login?redirect_to="+redirectURL)
		}
		return next(c)
	}
}
