package services

import (
	"github.com/lib/pq"
	"oauth2-provider/models"
	"oauth2-provider/storage"
	"testing"
)

func TestValidateAuthorizationRequest(t *testing.T) {
	store := storage.NewMemoryStorage()
	service := NewOAuthService(store)

	// Setup Client
	client := &models.Client{
		ClientID:     "client1",
		Secret:       "secret",
		RedirectURIs: pq.StringArray{"http://localhost/cb"},
		GrantTypes:   pq.StringArray{"authorization_code"},
	}
	store.StoreClient(client)

	req := &models.AuthorizationRequest{
		ClientID:            "client1",
		RedirectURI:         "http://localhost/cb",
		ResponseType:        "code",
		CodeChallenge:       "challenge",
		CodeChallengeMethod: "S256",
	}

	err := service.ValidateAuthorizationRequest(req)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	// Test Invalid Client
	req.ClientID = "invalid"
	err = service.ValidateAuthorizationRequest(req)
	if err == nil {
		t.Error("Expected error for invalid client")
	}

	// Reset ClientID
	req.ClientID = "client1"

	// Test Invalid Redirect URI
	req.RedirectURI = "http://evil.com"
	err = service.ValidateAuthorizationRequest(req)
	if err == nil {
		t.Error("Expected error for invalid redirect URI")
	}

	// Test Invalid Response Type
	req.RedirectURI = "http://localhost/cb"
	req.ResponseType = "token"
	err = service.ValidateAuthorizationRequest(req)
	if err == nil {
		t.Error("Expected error for invalid response type")
	}
}
