package services

import (
	"oauth2-provider/models"
	"oauth2-provider/storage"
	"testing"
)

func TestOAuthService_Flow(t *testing.T) {
	// Setup
	store := storage.NewMemoryStorage()
	service := NewOAuthService(store)

	clientID := "test-client"
	userID := uint(1)
	redirectURI := "/callback"
	codeChallenge := "challenge"
	codeChallengeMethod := "plain"

	// 1. Create Client
	client := &models.Client{
		ClientID:     clientID,
		RedirectURIs: []string{redirectURI},
	}
	store.StoreClient(client)

	// 2. Authorize (Generate Code)
	code, err := service.GenerateAuthorizationCode(clientID, userID, redirectURI, codeChallenge, codeChallengeMethod)
	if err != nil {
		t.Fatalf("GenerateAuthorizationCode failed: %v", err)
	}
	if code == "" {
		t.Fatal("Expected authorization code, got empty string")
	}

	// 3. Exchange Token
	tokenReq := &models.TokenRequest{
		GrantType:    "authorization_code",
		Code:         code,
		RedirectURI:  redirectURI,
		ClientID:     clientID,
		CodeVerifier: "challenge", // For plain method, verifier == challenge
	}

	accessToken, refreshToken, err := service.ExchangeToken(tokenReq)
	if err != nil {
		t.Fatalf("ExchangeToken failed: %v", err)
	}
	if accessToken == "" {
		t.Fatal("Expected access token")
	}
	if refreshToken == "" {
		t.Fatal("Expected refresh token")
	}

	// 4. Test Invalid Redirect URI
	badReq := &models.TokenRequest{
		GrantType:    "authorization_code",
		Code:         code, // Note: Code is used now, so this test might fail if reuse is allowed or not checked correctly in test logic.
                            // But ExchangeToken consumes the code. So we need a new code.
		RedirectURI:  "/bad-uri",
		ClientID:     clientID,
		CodeVerifier: "challenge",
	}

	// Generate new code for failure test
	code2, _ := service.GenerateAuthorizationCode(clientID, userID, redirectURI, codeChallenge, codeChallengeMethod)
	badReq.Code = code2

	_, _, err = service.ExchangeToken(badReq)
	if err == nil {
		t.Fatal("Expected error for redirect_uri mismatch, got nil")
	}
	if err.Error() != "redirect_uri mismatch" {
		t.Fatalf("Expected 'redirect_uri mismatch', got '%v'", err)
	}
}
