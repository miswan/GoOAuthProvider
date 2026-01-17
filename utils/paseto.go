package utils

import (
	"fmt"
	"github.com/o1egl/paseto"
	"oauth2-provider/config"
	"time"
)

type PasetoClaims struct {
	Subject   string    `json:"sub"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}

func GeneratePaseto(userID uint, duration time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(duration)

	jsonToken := PasetoClaims{
		Subject:   fmt.Sprintf("%d", userID),
		IssuedAt:  now,
		ExpiresAt: exp,
	}

	// Paseto v2 Local requires a 32 byte key.
	// We need to ensure the key is the correct length.
	// For this implementation, we will pad or truncate the config key if necessary,
	// or better, we should probably update the config to ensure it provides a valid key.
	// But to avoid breaking changes in config right now, let's fix the key length here.
	key := getSymmetricKey()

	return paseto.NewV2().Encrypt(key, jsonToken, nil)
}

func ValidatePaseto(tokenString string) (*PasetoClaims, error) {
	var claims PasetoClaims
	key := getSymmetricKey()

	err := paseto.NewV2().Decrypt(tokenString, key, &claims, nil)
	if err != nil {
		return nil, err
	}

	if time.Now().After(claims.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

func getSymmetricKey() []byte {
	// Paseto v2 requires exactly 32 bytes for the key.
	key := []byte(config.JWTSecret)
	if len(key) == 32 {
		return key
	}

	// Pad with zeros or truncate
	newKey := make([]byte, 32)
	copy(newKey, key)
	return newKey
}
