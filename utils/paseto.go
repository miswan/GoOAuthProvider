package utils

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"oauth2-provider/config"
)

type TokenClaims struct {
	UserID string    `json:"sub"`
	Exp    time.Time `json:"exp"`
}

func GeneratePaseto(userID uint, duration time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(duration)

	jsonToken := paseto.JSONToken{
		Audience:   "oauth2-provider",
		Issuer:     "oauth2-provider",
		Jti:        fmt.Sprintf("%d", now.UnixNano()),
		Subject:    fmt.Sprintf("%d", userID),
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  now,
	}

	v2 := paseto.NewV2()
	token, err := v2.Encrypt([]byte(config.PasetoKey), jsonToken, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidatePaseto(tokenString string) (*paseto.JSONToken, error) {
	v2 := paseto.NewV2()
	var token paseto.JSONToken
	var footer string

	err := v2.Decrypt(tokenString, []byte(config.PasetoKey), &token, &footer)
	if err != nil {
		return nil, err
	}

	if time.Now().After(token.Expiration) {
		return nil, fmt.Errorf("token expired")
	}

	return &token, nil
}
