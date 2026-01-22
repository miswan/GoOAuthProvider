package services

import (
	"oauth2-provider/models"
	"oauth2-provider/storage"
	"testing"
)

func TestOAuthService_GenerateAuthorizationCode(t *testing.T) {
	store := storage.NewMemoryStorage()
	s := NewOAuthService(store)

	code, err := s.GenerateAuthorizationCode("client_id", 1, "challenge", "S256", "http://localhost/callback")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if code == "" {
		t.Fatal("expected code, got empty string")
	}

	authCode := store.GetAuthCode(code)
	if authCode == nil {
		t.Fatal("expected auth code to be stored")
	}
	if authCode.RedirectURI != "http://localhost/callback" {
		t.Errorf("expected redirect URI %s, got %s", "http://localhost/callback", authCode.RedirectURI)
	}
}

func TestOAuthService_ExchangeToken(t *testing.T) {
	store := storage.NewMemoryStorage()
	s := NewOAuthService(store)

	// Setup
	code := "auth_code"
	clientID := "client_1"
	userID := uint(1)
	redirectURI := "http://localhost/callback"

	// We need to manually store auth code because GenerateAuthorizationCode generates random code
	store.StoreAuthCodeWithPKCE(code, clientID, userID, "challenge", "plain", redirectURI)

	// Test Success
	req := &models.TokenRequest{
		GrantType:    "authorization_code",
		Code:         code,
		RedirectURI:  redirectURI,
		ClientID:     clientID,
		CodeVerifier: "challenge", // plain
	}

	accessToken, refreshToken, err := s.ExchangeToken(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if accessToken == "" || refreshToken == "" {
		t.Fatal("expected tokens")
	}

	// Test Used Code
	// GetAuthCode in MemoryStorage marks it as used.
	// So calling ExchangeToken again should fail because GetAuthCode will return nil (or handle used check)
	_, _, err = s.ExchangeToken(req)
	if err == nil {
		t.Fatal("expected error on reused code")
	}

	// Test Mismatched Redirect URI
	store.StoreAuthCodeWithPKCE("code_2", clientID, userID, "challenge", "plain", redirectURI)
	req2 := &models.TokenRequest{
		GrantType:    "authorization_code",
		Code:         "code_2",
		RedirectURI:  "http://wrong.com",
		ClientID:     clientID,
		CodeVerifier: "challenge",
	}
	_, _, err = s.ExchangeToken(req2)
	if err == nil {
		t.Fatal("expected error on redirect uri mismatch")
	} else if err.Error() != "redirect uri mismatch" {
		t.Errorf("expected 'redirect uri mismatch', got '%v'", err)
	}
}
