package storage

import (
	"errors"
	"oauth2-provider/models"
	"sync"
	"time"
)

type MemoryStorage struct {
	users         map[uint]*models.User
	clients       map[string]*models.Client
	authCodes     map[string]*models.AuthCode
	refreshTokens map[string]*models.RefreshToken
	mu            sync.RWMutex
	nextUserID    uint
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users:         make(map[uint]*models.User),
		clients:       make(map[string]*models.Client),
		authCodes:     make(map[string]*models.AuthCode),
		refreshTokens: make(map[string]*models.RefreshToken),
		nextUserID:    1,
	}
}

func (s *MemoryStorage) StoreUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if user.ID == 0 {
		user.ID = s.nextUserID
		s.nextUserID++
	}
	s.users[user.ID] = user
	return nil
}

func (s *MemoryStorage) GetUserByUsername(username string) *models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, user := range s.users {
		if user.Username == username {
			return user
		}
	}
	return nil
}

func (s *MemoryStorage) GetUser(id uint) *models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users[id]
}

func (s *MemoryStorage) GetClient(clientID string) *models.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clients[clientID]
}

func (s *MemoryStorage) StoreClient(client *models.Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client.ClientID] = client
	return nil
}

func (s *MemoryStorage) StoreAuthCode(code, clientID string, userID uint) error {
	return s.StoreAuthCodeWithPKCE(code, clientID, userID, "", "", "")
}

func (s *MemoryStorage) StoreAuthCodeWithPKCE(code, clientID string, userID uint, codeChallenge, codeChallengeMethod, redirectURI string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.authCodes[code] = &models.AuthCode{
		Code:                code,
		ClientID:            clientID,
		UserID:              userID,
		ExpiresAt:           time.Now().Add(10 * time.Minute),
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		RedirectURI:         redirectURI,
		Used:                false,
	}
	return nil
}

func (s *MemoryStorage) GetAuthCode(code string) *models.AuthCode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if auth, exists := s.authCodes[code]; exists && time.Now().Before(auth.ExpiresAt) && !auth.Used {
		// Return a copy to prevent modification without calling MarkAuthCodeUsed
		val := *auth
		return &val
	}
	return nil
}

func (s *MemoryStorage) MarkAuthCodeUsed(code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if auth, exists := s.authCodes[code]; exists {
		auth.Used = true
		return nil
	}
	return errors.New("auth code not found")
}

func (s *MemoryStorage) StoreRefreshToken(token string, userID uint, clientID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refreshTokens[token] = &models.RefreshToken{
		Token:    token,
		UserID:   userID,
		ClientID: clientID,
		ExpiresAt: time.Now().Add(24 * time.Hour * 30),
	}
	return nil
}

func (s *MemoryStorage) GetRefreshToken(token string) *models.RefreshToken {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if rt, exists := s.refreshTokens[token]; exists && time.Now().Before(rt.ExpiresAt) {
		return rt
	}
	return nil
}

func (s *MemoryStorage) DeleteRefreshToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.refreshTokens, token)
	return nil
}
