package utils

import (
	"fmt"
	"github.com/o1egl/paseto"
	"oauth2-provider/config"
	"time"
)

func GeneratePaseto(userID uint, duration time.Duration) (string, error) {
	jsonToken := paseto.JSONToken{
		Subject:   fmt.Sprintf("%d", userID),
		IssuedAt:  time.Now(),
		Expiration: time.Now().Add(duration),
	}
	// Add footer if needed, currently empty
	footer := ""

	// Encrypt
	return paseto.NewV2().Encrypt([]byte(config.PasetoKey), jsonToken, footer)
}

func ValidatePaseto(tokenString string) (*paseto.JSONToken, error) {
	var jsonToken paseto.JSONToken
	var footer string
	err := paseto.NewV2().Decrypt(tokenString, []byte(config.PasetoKey), &jsonToken, &footer)
	if err != nil {
		return nil, err
	}

	// Check expiration
	if time.Now().After(jsonToken.Expiration) {
		return nil, fmt.Errorf("token expired")
	}

	return &jsonToken, nil
}
