package utils

import (
	"fmt"
	"time"
	"oauth2-provider/config"
	"github.com/o1egl/paseto"
)

type PasetoClaims struct {
	Subject   string    `json:"sub"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}

func GeneratePaseto(userID uint, duration time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(duration)

	jsonToken := paseto.JSONToken{
		Subject:    fmt.Sprintf("%d", userID),
		IssuedAt:   now,
		Expiration: exp,
	}

	return paseto.NewV2().Encrypt([]byte(config.JWTSecret), jsonToken, nil)
}

func ValidatePaseto(tokenString string) (*PasetoClaims, error) {
	var jsonToken paseto.JSONToken
	var footer string
	err := paseto.NewV2().Decrypt(tokenString, []byte(config.JWTSecret), &jsonToken, &footer)
	if err != nil {
		return nil, err
	}

    if err := jsonToken.Validate(); err != nil {
        return nil, err
    }

	return &PasetoClaims{
		Subject:   jsonToken.Subject,
		IssuedAt:  jsonToken.IssuedAt,
		ExpiresAt: jsonToken.Expiration,
	}, nil
}
