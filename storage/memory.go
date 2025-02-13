package storage

import (
	"oauth2-provider/models"
	"sync"
	"time"
)

type AuthCode struct {
	Code               string
	ClientID           string
	UserID             string
	ExpiresAt         time.Time
	CodeChallenge     string
	CodeChallengeMethod string
}

type MemoryStorage struct {
	users        map[string]*models.User
	clients      map[string]*models.Client
	authCodes    map[string]*AuthCode
	refreshTokens map[string]string
	mu           sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users:         make(map[string]*models.User),
		clients:       make(map[string]*models.Client),
		authCodes:     make(map[string]*AuthCode),
		refreshTokens: make(map[string]string),
	}
}

func (s *MemoryStorage) StoreUser(user *models.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
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

func (s *MemoryStorage) GetClient(clientID string) *models.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clients[clientID]
}

// Menambahkan method StoreClient yang hilang
func (s *MemoryStorage) StoreClient(client *models.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client.ID] = client
}

func (s *MemoryStorage) StoreAuthCode(code, clientID, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.authCodes[code] = &AuthCode{
		Code:      code,
		ClientID:  clientID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
}

func (s *MemoryStorage) GetAuthCode(code string) *AuthCode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if auth, exists := s.authCodes[code]; exists && time.Now().Before(auth.ExpiresAt) {
		return auth
	}
	return nil
}

func (s *MemoryStorage) StoreRefreshToken(token, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refreshTokens[token] = userID
}

// Menambahkan method untuk mendukung PKCE
func (s *MemoryStorage) StoreAuthCodeWithPKCE(code, clientID, userID, codeChallenge, codeChallengeMethod string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.authCodes[code] = &AuthCode{
		Code:               code,
		ClientID:           clientID,
		UserID:            userID,
		ExpiresAt:         time.Now().Add(10 * time.Minute),
		CodeChallenge:     codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
}

// Menambahkan method GetRefreshToken
func (s *MemoryStorage) GetRefreshToken(token string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.refreshTokens[token]
}