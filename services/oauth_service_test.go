package services_test

import (
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/storage"
	"testing"
)

func TestGenerateAuthorizationCode(t *testing.T) {
	store := storage.NewMemoryStorage()
	service := services.NewOAuthService(store)

	code, err := service.GenerateAuthorizationCode("client_id", 1, "challenge", "S256")
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}
	if code == "" {
		t.Fatal("Code is empty")
	}

	authCode := store.GetAuthCode(code)
	if authCode == nil {
		t.Fatal("Auth code not stored")
	}
}

func TestValidateAuthorizationRequest(t *testing.T) {
	store := storage.NewMemoryStorage()
	service := services.NewOAuthService(store)

	client := &models.Client{
		ClientID:     "test_client",
		RedirectURIs: []string{"http://localhost/callback"},
		GrantTypes:   []string{"authorization_code"},
	}
	store.StoreClient(client)

	req := &models.AuthorizationRequest{
		ClientID:            "test_client",
		RedirectURI:         "http://localhost/callback",
		ResponseType:        "code",
		CodeChallenge:       "challenge",
		CodeChallengeMethod: "S256",
	}

	if err := service.ValidateAuthorizationRequest(req); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	req.RedirectURI = "http://bad.com"
	if err := service.ValidateAuthorizationRequest(req); err == nil {
		t.Fatal("Validation should fail for bad redirect URI")
	}
}
