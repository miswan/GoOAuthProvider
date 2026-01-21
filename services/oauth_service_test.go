package services

import (
	"oauth2-provider/models"
	"oauth2-provider/storage"
	"oauth2-provider/utils"
	"testing"
)

func TestOAuthFlow(t *testing.T) {
	// Initialize memory storage
	store := storage.NewMemoryStorage()
	oauthService := NewOAuthService(store)
	clientService := NewClientService(store)
	userService := NewUserService(store)

	// Register User
	userReg := &models.UserRegister{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}
	if err := userService.Register(userReg); err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Login User to get ID
	user, err := userService.Login(&models.UserLogin{Username: "testuser", Password: "password123"})
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}

	// Register Client
	clientReg := &models.ClientRegistration{
		RedirectURIs: []string{"http://localhost:8080/callback"},
	}
	client, err := clientService.RegisterClient(clientReg)
	if err != nil {
		t.Fatalf("Failed to register client: %v", err)
	}

	// Generate Code Challenge
	codeVerifier := utils.GenerateRandomString(43)
	// For plain method, challenge is same as verifier
	codeChallenge := codeVerifier

	// Authorize (Validate Request)
	authReq := &models.AuthorizationRequest{
		ClientID:            client.ClientID,
		RedirectURI:         "http://localhost:8080/callback",
		ResponseType:        "code",
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: "plain",
	}
	if err := oauthService.ValidateAuthorizationRequest(authReq); err != nil {
		t.Fatalf("ValidateAuthorizationRequest failed: %v", err)
	}

	// Generate Code (Simulation of OAuthHandler logic)
	code, err := oauthService.GenerateAuthorizationCode(client.ClientID, user.ID, authReq.RedirectURI, codeChallenge, "plain")
	if err != nil {
		t.Fatalf("GenerateAuthorizationCode failed: %v", err)
	}

	// Exchange Token
	tokenReq := &models.TokenRequest{
		GrantType:    "authorization_code",
		Code:         code,
		RedirectURI:  "http://localhost:8080/callback",
		ClientID:     client.ClientID,
		CodeVerifier: codeVerifier,
	}

	accessToken, refreshToken, err := oauthService.ExchangeToken(tokenReq)
	if err != nil {
		t.Fatalf("ExchangeToken failed: %v", err)
	}

	if accessToken == "" || refreshToken == "" {
		t.Fatal("Tokens should not be empty")
	}

	// Refresh Token
	refreshReq := &models.TokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}

	newAccess, newRefresh, err := oauthService.ExchangeToken(refreshReq)
	if err != nil {
		t.Fatalf("RefreshToken grant failed: %v", err)
	}

	if newAccess == "" || newRefresh == "" {
		t.Fatal("New tokens should not be empty")
	}
}
