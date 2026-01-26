package middleware

import (
	"github.com/labstack/echo/v4"
	"oauth2-provider/utils"
	"strings"
)

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := ""

		// Check header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// Check cookie if header missing
		if tokenString == "" {
			cookie, err := c.Cookie("auth_token")
			if err == nil {
				tokenString = cookie.Value
			}
		}

		if tokenString == "" {
			return echo.ErrUnauthorized
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return echo.ErrUnauthorized
		}

		c.Set("user_id", claims.Subject)
		return next(c)
	}
}
