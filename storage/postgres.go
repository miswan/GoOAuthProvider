package storage

import (
	"gorm.io/gorm"
	"oauth2-provider/models"
	"oauth2-provider/utils"
	"time"
)

type PostgresStorage struct {
	db *gorm.DB
}

func NewPostgresStorage(db *gorm.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) StoreUser(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *PostgresStorage) GetUserByUsername(username string) *models.User {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil
	}
	return &user
}

func (s *PostgresStorage) StoreClient(client *models.Client) error {
	// Generate client credentials
	client.ClientID = utils.GenerateRandomString(24)
	client.Secret = utils.GenerateRandomString(32)

	// Store client with arrays
	return s.db.Create(client).Error
}

func (s *PostgresStorage) GetClient(clientID string) *models.Client {
	var client models.Client
	if err := s.db.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		return nil
	}
	return &client
}

func (s *PostgresStorage) StoreAuthCode(code, clientID string, userID uint) error {
	authCode := &models.AuthCode{
		Code:      code,
		ClientID:  clientID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	return s.db.Create(authCode).Error
}

func (s *PostgresStorage) StoreAuthCodeWithPKCE(code, clientID string, userID uint, codeChallenge, codeChallengeMethod string) error {
	authCode := &models.AuthCode{
		Code:                code,
		ClientID:            clientID,
		UserID:             userID,
		ExpiresAt:          time.Now().Add(10 * time.Minute),
		CodeChallenge:      codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
	return s.db.Create(authCode).Error
}

func (s *PostgresStorage) GetAuthCode(code string) *models.AuthCode {
	var authCode models.AuthCode
	if err := s.db.Where("code = ? AND expires_at > ? AND used = ?", code, time.Now(), false).First(&authCode).Error; err != nil {
		return nil
	}

	// Mark the auth code as used
	s.db.Model(&authCode).Update("used", true)

	return &authCode
}

func (s *PostgresStorage) StoreRefreshToken(token string, userID uint, clientID string) error {
	refreshToken := &models.RefreshToken{
		Token:     token,
		UserID:    userID,
		ClientID:  clientID,
		ExpiresAt: time.Now().Add(24 * time.Hour * 30), // 30 days
	}
	return s.db.Create(refreshToken).Error
}

func (s *PostgresStorage) GetRefreshToken(token string) *models.RefreshToken {
	var refreshToken models.RefreshToken
	if err := s.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken).Error; err != nil {
		return nil
	}
	return &refreshToken
}

func (s *PostgresStorage) DeleteRefreshToken(token string) error {
	return s.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}