package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"oauth2-provider/models"
	"testing"

	"github.com/labstack/echo/v4"
)


func TestOAuthHandler_Authorize_Redirects(t *testing.T) {
	// Mock implementation
	mockService := &MockOAuthService{}
	h := NewOAuthHandler(mockService) // This will fail if I don't change the constructor

	// Setup Echo
	e := echo.New()
	q := make(url.Values)
	q.Set("client_id", "client123")
	q.Set("redirect_uri", "http://example.com/cb")
	q.Set("response_type", "code")
	q.Set("state", "xyz")
	q.Set("code_challenge", "challenge")
	q.Set("code_challenge_method", "plain")

	req := httptest.NewRequest(http.MethodGet, "/authorize?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	if err := h.Authorize(c); err != nil {
		t.Errorf("Authorize returned error: %v", err)
	}

	// Assert
	if rec.Code != http.StatusFound {
		t.Errorf("Expected status 302, got %d", rec.Code)
	}

	loc := rec.Header().Get("Location")
	parsedLoc, _ := url.Parse(loc)
	if parsedLoc.Scheme != "http" || parsedLoc.Host != "example.com" || parsedLoc.Path != "/cb" {
		t.Errorf("Redirect location mismatch: %s", loc)
	}
	if parsedLoc.Query().Get("code") != "mock_code" {
		t.Errorf("Expected code 'mock_code', got %s", parsedLoc.Query().Get("code"))
	}
	if parsedLoc.Query().Get("state") != "xyz" {
		t.Errorf("Expected state 'xyz', got %s", parsedLoc.Query().Get("state"))
	}
}

type MockOAuthService struct{}

func (m *MockOAuthService) ValidateAuthorizationRequest(req *models.AuthorizationRequest) error {
	return nil
}

func (m *MockOAuthService) GenerateAuthorizationCode(clientID string, userID uint, codeChallenge, codeChallengeMethod string) (string, error) {
	return "mock_code", nil
}

func (m *MockOAuthService) ExchangeToken(req *models.TokenRequest) (string, string, error) {
	return "access", "refresh", nil
}
