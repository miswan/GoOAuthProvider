package utils

import (
	"fmt"
	"oauth2-provider/config"
	"time"

	"github.com/o1egl/paseto"
)

// GeneratePaseto creates a new Paseto v2 Local token
func GeneratePaseto(userID uint, duration time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(duration)

	jsonToken := paseto.JSONToken{
		Audience:   "oauth2-provider",
		Issuer:     "oauth2-provider",
		Jti:        GenerateRandomString(16),
		Subject:    fmt.Sprintf("%d", userID),
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  now,
	}

	// Encrypt
	token, err := paseto.NewV2().Encrypt(config.PasetoKey, jsonToken, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidatePaseto checks the validity of a Paseto token and returns the JSONToken
func ValidatePaseto(tokenString string) (*paseto.JSONToken, error) {
	var jsonToken paseto.JSONToken
	var footer string

	err := paseto.NewV2().Decrypt(tokenString, config.PasetoKey, &jsonToken, &footer)
	if err != nil {
		return nil, err
	}

	// Validate expiration
	if time.Now().After(jsonToken.Expiration) {
		return nil, fmt.Errorf("token expired")
	}

	return &jsonToken, nil
}
