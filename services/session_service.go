package services

import (
	"oauth2-provider/utils"
	"time"
)

func GenerateSessionToken(userID uint) (string, error) {
	return utils.GenerateJWT(userID, 24*time.Hour)
}
