package utils

import (
	"fmt"
	"oauth2-provider/config"
	"time"

	"github.com/o1egl/paseto"
)

type PasetoClaims struct {
	Subject   string    `json:"sub"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}

func GeneratePaseto(userID uint, duration time.Duration) (string, error) {
	pst := paseto.NewV2()
	key := []byte(config.JWTSecret)

	claims := PasetoClaims{
		Subject:   fmt.Sprintf("%d", userID),
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	return pst.Encrypt(key, claims, nil)
}

func ValidatePaseto(tokenString string) (*PasetoClaims, error) {
	pst := paseto.NewV2()
	key := []byte(config.JWTSecret)

	var claims PasetoClaims
	err := pst.Decrypt(tokenString, key, &claims, nil)
	if err != nil {
		return nil, err
	}

	if time.Now().After(claims.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}
