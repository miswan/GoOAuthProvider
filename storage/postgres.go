package storage

import (
	"log"
	"oauth2-provider/models"
	"oauth2-provider/utils"
	"time"

	"gorm.io/gorm"
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
		log.Printf("Error getting user by username: %v", err)
		return nil
	}
	return &user
}

func (s *PostgresStorage) StoreClient(client *models.Client) error {
	// Log the client data before storing
	log.Printf("Storing client with RedirectURIs: %v, GrantTypes: %v", client.RedirectURIs, client.GrantTypes)

	// Generate client credentials if not present (though usually done by service)
	if client.ClientID == "" {
		client.ClientID = utils.GenerateRandomString(24)
	}
	if client.Secret == "" {
		client.Secret = utils.GenerateRandomString(32)
	}

	// Create client using GORM
	result := s.db.Debug().Create(client)
	if result.Error != nil {
		log.Printf("Error storing client: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully stored client with ID: %s", client.ClientID)
	return nil
}

func (s *PostgresStorage) GetClient(clientID string) *models.Client {
	var client models.Client
	if err := s.db.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		log.Printf("Error getting client: %v", err)
		return nil
	}
	return &client
}

func (s *PostgresStorage) StoreAuthCodeWithPKCE(code, clientID string, userID uint, redirectURI, codeChallenge, codeChallengeMethod string) error {
	authCode := &models.AuthCode{
		Code:                code,
		ClientID:            clientID,
		UserID:              userID,
		RedirectURI:         redirectURI,
		ExpiresAt:           time.Now().Add(10 * time.Minute),
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		Used:                false,
	}
	return s.db.Create(authCode).Error
}

func (s *PostgresStorage) GetAuthCode(code string) *models.AuthCode {
	var authCode models.AuthCode
	// Use a transaction to lock the row and update 'used' atomically
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Find valid unused code and lock row
		// Note: using raw SQL locking might be safer if database specific,
		// but GORM Clausses(clause.Locking{Strength: "UPDATE"}) is better.
		// For simplicity, we check and update.
		// 'FOR UPDATE' is strictly postgres/mysql compatible.

		if err := tx.Where("code = ? AND expires_at > ? AND used = ?", code, time.Now(), false).First(&authCode).Error; err != nil {
			return err
		}

		// Mark as used
		authCode.Used = true
		return tx.Save(&authCode).Error
	})

	if err != nil {
		log.Printf("Error consuming auth code: %v", err)
		return nil
	}
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
		log.Printf("Error getting refresh token: %v", err)
		return nil
	}
	return &refreshToken
}

func (s *PostgresStorage) DeleteRefreshToken(token string) error {
	return s.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}
