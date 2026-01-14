package services

import (
	"oauth2-provider/utils"
	"time"
)

// GenerateSessionToken wraps the utility function to allow for consistent use in services
func GenerateSessionToken(userID uint) (string, error) {
	return utils.GenerateJWT(userID, 24*time.Hour)
}
