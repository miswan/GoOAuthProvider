package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"oauth2-provider/utils"
	"net/url"
)

func SessionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			// Redirect to login if cookie is missing
			// Preserve current URL as return destination
			currentURL := c.Request().URL.String()
			return c.Redirect(http.StatusFound, "/login?redirect_to="+url.QueryEscape(currentURL))
		}

		// Validate token (assuming it's a JWT)
		claims, err := utils.ValidateJWT(cookie.Value)
		if err != nil {
			return c.Redirect(http.StatusFound, "/login")
		}

		// Store user ID in context
		c.Set("user_id", claims.Subject)
		return next(c)
	}
}
