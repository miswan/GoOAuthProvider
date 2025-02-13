package utils

import (
    "github.com/golang-jwt/jwt"
    "oauth2-provider/config"
    "time"
)

func GenerateJWT(userID string, duration time.Duration) (string, error) {
    claims := jwt.StandardClaims{
        Subject:   userID,
        IssuedAt:  time.Now().Unix(),
        ExpiresAt: time.Now().Add(duration).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.JWTSecret))
}

func ValidateJWT(tokenString string) (*jwt.StandardClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(config.JWTSecret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, jwt.ErrSignatureInvalid
}
