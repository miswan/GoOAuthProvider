package services

import (
	"oauth2-provider/models"
	"oauth2-provider/utils"
	"time"
)

func GenerateSessionToken(user *models.User) (string, error) {
	// Reusing JWT util for session token
	return utils.GenerateJWT(user.ID, 24*time.Hour)
}
