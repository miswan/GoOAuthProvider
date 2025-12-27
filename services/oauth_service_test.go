package services

import (
	"oauth2-provider/models"
	"testing"
)

// Mock Storage implementation
type MockStorage struct {
	clients     map[string]*models.Client
	authCodes   map[string]*models.AuthCode
	refreshTokens map[string]*models.RefreshToken
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		clients:     make(map[string]*models.Client),
		authCodes:   make(map[string]*models.AuthCode),
		refreshTokens: make(map[string]*models.RefreshToken),
	}
}

func (m *MockStorage) GetClient(clientID string) *models.Client {
	return m.clients[clientID]
}

func (m *MockStorage) StoreClient(client *models.Client) error {
	m.clients[client.ClientID] = client
	return nil
}

func (m *MockStorage) StoreAuthCodeWithPKCE(code, clientID string, userID uint, redirectURI, codeChallenge, codeChallengeMethod string) error {
	m.authCodes[code] = &models.AuthCode{
		Code:                code,
		ClientID:            clientID,
		UserID:             userID,
		RedirectURI:        redirectURI,
		CodeChallenge:      codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
	return nil
}

func (m *MockStorage) GetAuthCode(code string) *models.AuthCode {
	return m.authCodes[code]
}

func (m *MockStorage) StoreRefreshToken(token string, userID uint, clientID string) error {
	m.refreshTokens[token] = &models.RefreshToken{
		Token:    token,
		UserID:   userID,
		ClientID: clientID,
	}
	return nil
}

func (m *MockStorage) GetRefreshToken(token string) *models.RefreshToken {
	return m.refreshTokens[token]
}

func (m *MockStorage) DeleteRefreshToken(token string) error {
	delete(m.refreshTokens, token)
	return nil
}

// Unused methods for this test
func (m *MockStorage) StoreUser(user *models.User) error { return nil }
func (m *MockStorage) GetUserByUsername(username string) *models.User { return nil }
func (m *MockStorage) StoreAuthCode(code, clientID string, userID uint) error { return nil }

func TestExchangeToken_AuthorizationCode(t *testing.T) {
	mockStore := NewMockStorage()
	service := NewOAuthService(mockStore)

	// Setup Client
	clientID := "test-client"
	clientSecret := "test-secret"
	mockStore.StoreClient(&models.Client{
		ClientID: clientID,
		Secret:   clientSecret,
		RedirectURIs: []string{"http://localhost/callback"},
	})

	// Setup Auth Code
	authCode := "test-auth-code"
	verifier := "test-verifier"

	// Let's use plain for simplicity as the service supports it
	challenge := verifier
	method := "plain"

	mockStore.StoreAuthCodeWithPKCE(authCode, clientID, 1, "http://localhost/callback", challenge, method)

	req := &models.TokenRequest{
		GrantType:    "authorization_code",
		Code:         authCode,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  "http://localhost/callback",
		CodeVerifier: verifier,
	}

	accessToken, refreshToken, err := service.ExchangeToken(req)
	if err != nil {
		t.Fatalf("ExchangeToken failed: %v", err)
	}

	if accessToken == "" || refreshToken == "" {
		t.Error("Expected access token and refresh token")
	}
}

func TestExchangeToken_InvalidClientSecret(t *testing.T) {
	mockStore := NewMockStorage()
	service := NewOAuthService(mockStore)

	clientID := "test-client"
	clientSecret := "test-secret"
	mockStore.StoreClient(&models.Client{
		ClientID: clientID,
		Secret:   clientSecret,
	})

	req := &models.TokenRequest{
		GrantType:    "authorization_code",
		Code:         "code",
		ClientID:     clientID,
		ClientSecret: "wrong-secret",
	}

	_, _, err := service.ExchangeToken(req)
	if err == nil {
		t.Error("Expected error for invalid client secret, got nil")
	}
}

func TestExchangeToken_RedirectURIMismatch(t *testing.T) {
	mockStore := NewMockStorage()
	service := NewOAuthService(mockStore)

	clientID := "test-client"
	authCode := "code"
	mockStore.StoreClient(&models.Client{ClientID: clientID})

	// Store code with one URI
	mockStore.StoreAuthCodeWithPKCE(authCode, clientID, 1, "http://localhost/callback", "challenge", "plain")

	// Request with different URI
	req := &models.TokenRequest{
		GrantType:    "authorization_code",
		Code:         authCode,
		ClientID:     clientID,
		RedirectURI:  "http://localhost/other",
		CodeVerifier: "verifier",
	}

	_, _, err := service.ExchangeToken(req)
	if err == nil {
		t.Error("Expected error for redirect URI mismatch, got nil")
	} else if err.Error() != "redirect_uri mismatch" {
		t.Errorf("Expected 'redirect_uri mismatch', got '%v'", err)
	}
}
