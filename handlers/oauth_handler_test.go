package handlers

import (
	"net/http"
	"net/http/httptest"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/storage"
	"oauth2-provider/utils"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*OAuthHandler, *storage.MemoryStorage, *services.OAuthService, *services.UserService) {
	store := storage.NewMemoryStorage()
	// Setup initial data
	store.StoreClient(&models.Client{
		ClientID:     "client123",
		Secret:       "secret",
		RedirectURIs: []string{"http://localhost:3000/callback"},
	})
	store.StoreUser(&models.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
		Password: "password",
	})

	oauthService := services.NewOAuthService(store)
	userService := services.NewUserService(store)
	handler := NewOAuthHandler(oauthService, userService)

	return handler, store, oauthService, userService
}

func TestAuthorizeRedirectsToLogin(t *testing.T) {
	handler, _, _, _ := setupTest(t)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/authorize?client_id=client123&redirect_uri=http://localhost:3000/callback&response_type=code&code_challenge=foo&code_challenge_method=plain", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.Authorize(c); err != nil {
		t.Errorf("Authorize returned error: %v", err)
	}

	if rec.Code != http.StatusFound {
		t.Errorf("Expected status 302, got %d", rec.Code)
	}

	location := rec.Header().Get("Location")
	if !strings.Contains(location, "/login?return_to=") {
		t.Errorf("Expected redirect to login, got %s", location)
	}
}

func TestAuthorizeSuccess(t *testing.T) {
	handler, _, _, _ := setupTest(t)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/authorize?client_id=client123&redirect_uri=http://localhost:3000/callback&response_type=code&state=xyz&code_challenge=foo&code_challenge_method=plain", nil)

	// Set Auth Cookie
	token, _ := utils.GenerateJWT(1, time.Hour)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.Authorize(c); err != nil {
		t.Errorf("Authorize returned error: %v", err)
	}

	if rec.Code != http.StatusFound {
		t.Errorf("Expected status 302, got %d", rec.Code)
	}

	location := rec.Header().Get("Location")
	if !strings.HasPrefix(location, "http://localhost:3000/callback") {
		t.Errorf("Expected redirect to callback, got %s", location)
	}
	if !strings.Contains(location, "code=") {
		t.Errorf("Expected code in redirect, got %s", location)
	}
	if !strings.Contains(location, "state=xyz") {
		t.Errorf("Expected state in redirect, got %s", location)
	}
}
