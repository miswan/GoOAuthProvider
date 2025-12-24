package middleware

import (
    "github.com/labstack/echo/v4"
    "oauth2-provider/utils"
    "net/http"
    "net/url"
)

// SessionAuthMiddleware checks for a session cookie and sets user_id in context if valid.
// It does NOT redirect or error out if session is missing; that is handled by the handler.
// This allows the handler to decide whether to redirect to login or proceed.
// Wait, the handler needs to know if user is authenticated.
// If I use this middleware, I can check c.Get("user_id").
func SessionAuth(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        cookie, err := c.Cookie("session_token")
        if err != nil {
            // No cookie, proceed without user_id
            return next(c)
        }

        claims, err := utils.ValidateJWT(cookie.Value)
        if err != nil {
            // Invalid cookie, proceed without user_id (maybe clear cookie?)
            return next(c)
        }

        c.Set("user_id", claims.Subject)
        return next(c)
    }
}

// RequireSessionAuth works like SessionAuth but redirects to Login if session is missing.
// This is useful for endpoints that MUST have a user (like /authorize).
// However, /authorize needs to preserve the query params (return_to).
func RequireSessionAuth(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        cookie, err := c.Cookie("session_token")
        if err != nil || cookie.Value == "" {
             return redirectToLogin(c)
        }

        claims, err := utils.ValidateJWT(cookie.Value)
        if err != nil {
            return redirectToLogin(c)
        }

        c.Set("user_id", claims.Subject)
        return next(c)
    }
}

func redirectToLogin(c echo.Context) error {
    // Current URL as return_to
    currentURL := c.Request().URL.String()
    return c.Redirect(http.StatusFound, "/login?return_to="+url.QueryEscape(currentURL))
}
