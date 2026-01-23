package middleware

import (
    "github.com/labstack/echo/v4"
    "oauth2-provider/utils"
    "strconv"
    "strings"
)

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token string
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		if token == "" {
			cookie, err := c.Cookie("auth_token")
			if err == nil {
				token = cookie.Value
			}
		}

		if token == "" {
			return echo.ErrUnauthorized
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			return echo.ErrUnauthorized
		}

		uid, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			return echo.ErrUnauthorized
		}

		c.Set("user_id", uint(uid))
		return next(c)
	}
}
