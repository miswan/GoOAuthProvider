package storage

import "oauth2-provider/models"

type Storage interface {
	StoreUser(user *models.User) error
	GetUserByUsername(username string) *models.User

	StoreClient(client *models.Client) error
	GetClient(clientID string) *models.Client

	StoreAuthCode(code, clientID string, userID uint, redirectURI string) error
	StoreAuthCodeWithPKCE(code, clientID string, userID uint, redirectURI, codeChallenge, codeChallengeMethod string) error
	GetAuthCode(code string) *models.AuthCode

	StoreRefreshToken(token string, userID uint, clientID string) error
	GetRefreshToken(token string) *models.RefreshToken
	DeleteRefreshToken(token string) error
}
