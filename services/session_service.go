package services

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"oauth2-provider/config"
	"strconv"
)

type SessionClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func ValidateSessionToken(tokenString string) (*SessionClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		// Re-map to our internal structure (assuming Subject holds the UserID)
		uid, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user id in token")
		}
		return &SessionClaims{
			UserID: uint(uid),
			StandardClaims: *claims,
		}, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
