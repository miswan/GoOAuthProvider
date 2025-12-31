package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"oauth2-provider/utils"
	"strconv"
)

// SessionAuthMiddleware checks for a valid session_token cookie.
// If valid, it sets the user_id in the context.
// It does NOT block execution; it just sets context. Handlers must check validity.
// Alternatively, we could have a RequireSession middleware.
func SessionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			// No session cookie, proceed without user context
			return next(c)
		}

		claims, err := utils.ValidateJWT(cookie.Value)
		if err != nil {
			// Invalid token, proceed without user context
			return next(c)
		}

		// Set user_id in context (claims.Subject is usually string, convert if needed or keep as string)
		// Assuming Subject is user ID string
		c.Set("user_id", claims.Subject)

		// Also parse to uint for easier usage if consistent with app
		if uid, err := strconv.ParseUint(claims.Subject, 10, 64); err == nil {
			c.Set("user_id_uint", uint(uid))
		}

		return next(c)
	}
}

// RequireSessionAuth enforces that a user is logged in.
// If not, it redirects to /login with return_to parameter.
func RequireSessionAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// This relies on SessionAuthMiddleware running first
		userID := c.Get("user_id")
		if userID == nil {
			// Construct return URL
			req := c.Request()
			returnTo := req.URL.String()
			// Encode returnTo? No, Echo redirect handles it, but we should probably URL encode it
			// if we were constructing it manually. Here we just pass it.
			return c.Redirect(http.StatusFound, "/login?redirect_to="+returnTo)
		}
		return next(c)
	}
}
