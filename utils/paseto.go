package utils

import (
	"fmt"
	"github.com/o1egl/paseto"
	"oauth2-provider/config"
	"time"
)

type TokenClaims struct {
	UserID    uint      `json:"user_id"`
	ExpiresAt time.Time `json:"exp"`
}

func GeneratePaseto(userID uint, duration time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(duration)

	claims := TokenClaims{
		UserID:    userID,
		ExpiresAt: exp,
	}

	footer := "oauth2-provider"
	return paseto.NewV2().Encrypt([]byte(config.PasetoKey), claims, footer)
}

func ValidatePaseto(tokenString string) (*TokenClaims, error) {
	var claims TokenClaims
	var footer string

	err := paseto.NewV2().Decrypt(tokenString, []byte(config.PasetoKey), &claims, &footer)
	if err != nil {
		return nil, err
	}

	if time.Now().After(claims.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}
