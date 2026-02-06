package storage

import "oauth2-provider/models"

type Storage interface {
	// User operations
	GetUser(id uint) *models.User
	GetUserByUsername(username string) *models.User
	StoreUser(user *models.User) error

	// Client operations
	GetClient(clientID string) *models.Client
	StoreClient(client *models.Client) error

	// AuthCode operations
	StoreAuthCode(code string, clientID string, userID uint, redirectURI string, codeChallenge, codeChallengeMethod string) error
	GetAuthCode(code string) *models.AuthCode
	MarkAuthCodeUsed(code string) error

	// RefreshToken operations
	StoreRefreshToken(token string, userID uint, clientID string) error
	GetRefreshToken(token string) *models.RefreshToken
	DeleteRefreshToken(token string) error
}
