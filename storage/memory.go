package storage

import (
	"oauth2-provider/models"
	"strconv"
	"sync"
	"time"
)

type AuthCode struct {
	Code               string
	ClientID           string
	UserID             uint
	ExpiresAt          time.Time
	CodeChallenge      string
	CodeChallengeMethod string
}

type MemoryStorage struct {
	users         map[uint]*models.User
	clients       map[string]*models.Client
	authCodes     map[string]*AuthCode
	refreshTokens map[string]string
	mu            sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users:         make(map[uint]*models.User),
		clients:       make(map[string]*models.Client),
		authCodes:     make(map[string]*AuthCode),
		refreshTokens: make(map[string]string),
	}
}

func (s *MemoryStorage) StoreUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
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
	s.mu.Lock()
	defer s.mu.Unlock()
	s.authCodes[code] = &AuthCode{
		Code:      code,
		ClientID:  clientID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	return nil
}

func (s *MemoryStorage) GetAuthCode(code string) *AuthCode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if auth, exists := s.authCodes[code]; exists && time.Now().Before(auth.ExpiresAt) {
		return auth
	}
	return nil
}

func (s *MemoryStorage) StoreRefreshToken(token string, userID uint, clientID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refreshTokens[token] = strconv.FormatUint(uint64(userID), 10)
	return nil
}

func (s *MemoryStorage) StoreAuthCodeWithPKCE(code, clientID string, userID uint, codeChallenge, codeChallengeMethod string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.authCodes[code] = &AuthCode{
		Code:               code,
		ClientID:           clientID,
		UserID:             userID,
		ExpiresAt:          time.Now().Add(10 * time.Minute),
		CodeChallenge:      codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
	return nil
}

func (s *MemoryStorage) GetRefreshToken(token string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.refreshTokens[token]
}

func (s *MemoryStorage) DeleteRefreshToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.refreshTokens, token)
	return nil
}