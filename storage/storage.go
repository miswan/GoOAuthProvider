package storage

import "oauth2-provider/models"

type Storage interface {
	// User operations
	StoreUser(user *models.User) error
	GetUserByUsername(username string) *models.User
	GetUser(id uint) *models.User

	// Client operations
	StoreClient(client *models.Client) error
	GetClient(clientID string) *models.Client

	// Auth Code operations
	StoreAuthCodeWithPKCE(code, clientID string, userID uint, redirectURI, codeChallenge, codeChallengeMethod string) error
	GetAuthCode(code string) *models.AuthCode
	MarkAuthCodeUsed(code string) error

	// Refresh Token operations
	StoreRefreshToken(token string, userID uint, clientID string) error
	GetRefreshToken(token string) *models.RefreshToken
	DeleteRefreshToken(token string) error
}
