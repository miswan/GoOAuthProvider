package services

import (
	"github.com/golang-jwt/jwt"
	"oauth2-provider/utils"
)

// ValidateSessionToken wraps the utility function
func ValidateSessionToken(tokenString string) (*jwt.StandardClaims, error) {
	return utils.ValidateJWT(tokenString)
}
