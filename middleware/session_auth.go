package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"oauth2-provider/utils"
	"time"
)

// SessionAuthMiddleware checks for a valid session token in the cookie
func SessionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			return next(c) // No session, let the handler decide (e.g. redirect to login)
		}

		claims, err := utils.ValidateJWT(cookie.Value)
		if err != nil {
			return next(c) // Invalid session
		}

		// Set user_id in context
		c.Set("user_id", claims.Subject)
		return next(c)
	}
}

// SetSessionCookie sets a secure session cookie
func SetSessionCookie(c echo.Context, token string) {
	cookie := new(http.Cookie)
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true
	// cookie.Secure = true // Enable in production
	// cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)
}
