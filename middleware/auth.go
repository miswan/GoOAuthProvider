package middleware

import (
	"github.com/labstack/echo/v4"
	"oauth2-provider/utils"
	"strings"
	"net/http"
)

func SessionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			// No session cookie, continue without user_id
			return next(c)
		}

		token := cookie.Value
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			// Invalid token, maybe clear cookie?
			return next(c)
		}

		c.Set("user_id", claims.Subject)
		return next(c)
	}
}

func RequireSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("user_id")
		if userID == nil {
			// Redirect to login
			redirectURL := "/login"
			if c.Request().Method == http.MethodGet {
				redirectURL += "?redirect_to=" + c.Request().URL.String()
			}
			return c.Redirect(http.StatusFound, redirectURL)
		}
		return next(c)
	}
}

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        authHeader := c.Request().Header.Get("Authorization")
        if authHeader == "" {
            return echo.ErrUnauthorized
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return echo.ErrUnauthorized
        }

        token := parts[1]
        claims, err := utils.ValidateJWT(token)
        if err != nil {
            return echo.ErrUnauthorized
        }

        c.Set("user_id", claims.Subject)
        return next(c)
    }
}
