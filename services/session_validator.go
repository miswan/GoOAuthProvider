package services

import (
	"github.com/golang-jwt/jwt"
	"oauth2-provider/utils"
)

func ValidateSessionToken(token string) (*jwt.StandardClaims, error) {
	return utils.ValidateJWT(token)
}
