package middleware

import (
    "github.com/labstack/echo/v4"
    "oauth2-provider/utils"
    "strings"
)

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
