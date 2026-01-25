package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"oauth2-provider/handlers"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/storage"
	"oauth2-provider/utils"
)

func TestOAuthFlow(t *testing.T) {
	// Setup
	e := echo.New()
	store := storage.NewMemoryStorage()
	oauthService := services.NewOAuthService(store)
	userService := services.NewUserService(store)
	clientService := services.NewClientService(store)

	oauthHandler := handlers.NewOAuthHandler(oauthService, userService)

	// 1. Register Client
    clientReg := &models.ClientRegistration{
        RedirectURIs: []string{"http://localhost:3000/callback"},
    }
    client, err := clientService.RegisterClient(clientReg)
    if err != nil {
        t.Fatalf("Failed to register client: %v", err)
    }

	// 2. Register User
    userReg := &models.UserRegister{
        Username: "testuser",
        Password: "password123",
        Email:    "test@example.com",
    }
    if err := userService.Register(userReg); err != nil {
        t.Fatalf("Failed to register user: %v", err)
    }
    user := store.GetUserByUsername("testuser")

    // 3. Login (Generate Cookie)
    token, _ := utils.GenerateJWT(user.ID, time.Hour)
    cookie := &http.Cookie{
        Name:  "auth_token",
        Value: token,
    }

    // 4. Test Authorize (No Cookie) -> Redirect Login
    q := make(url.Values)
    q.Set("client_id", client.ClientID)
    q.Set("redirect_uri", "http://localhost:3000/callback")
    q.Set("response_type", "code")
    q.Set("state", "xyz")
    q.Set("code_challenge", "plain_challenge")
    q.Set("code_challenge_method", "plain")

    req := httptest.NewRequest(http.MethodGet, "/authorize?"+q.Encode(), nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if err := oauthHandler.Authorize(c); err != nil {
        t.Errorf("Authorize failed: %v", err)
    }
    if rec.Code != http.StatusFound {
        t.Errorf("Expected 302, got %d", rec.Code)
    }
    if !strings.Contains(rec.Header().Get("Location"), "/login") {
        t.Errorf("Expected redirect to login, got %s", rec.Header().Get("Location"))
    }

    // 5. Test Authorize (With Cookie) -> Redirect Client
    req = httptest.NewRequest(http.MethodGet, "/authorize?"+q.Encode(), nil)
    req.AddCookie(cookie)
    rec = httptest.NewRecorder()
    c = e.NewContext(req, rec)

    if err := oauthHandler.Authorize(c); err != nil {
        t.Errorf("Authorize failed: %v", err)
    }
    if rec.Code != http.StatusFound {
        t.Errorf("Expected 302, got %d", rec.Code)
    }
    loc, _ := url.Parse(rec.Header().Get("Location"))
    if !strings.HasPrefix(loc.String(), "http://localhost:3000/callback") {
         t.Errorf("Expected redirect to callback, got %s", loc.String())
    }
    code := loc.Query().Get("code")
    if code == "" {
        t.Fatal("Authorization code missing")
    }

    // 6. Test Token Exchange
    form := make(url.Values)
    form.Set("grant_type", "authorization_code")
    form.Set("code", code)
    form.Set("redirect_uri", "http://localhost:3000/callback")
    form.Set("client_id", client.ClientID)
    form.Set("client_secret", client.Secret)
    form.Set("code_verifier", "plain_challenge")

    req = httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(form.Encode()))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
    rec = httptest.NewRecorder()
    c = e.NewContext(req, rec)

    if err := oauthHandler.Token(c); err != nil {
        t.Errorf("Token failed: %v", err)
    }
    if rec.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d. Body: %s", rec.Code, rec.Body.String())
    }

    var tokenResp map[string]interface{}
    json.Unmarshal(rec.Body.Bytes(), &tokenResp)
    accessToken, ok := tokenResp["access_token"].(string)
    if !ok || accessToken == "" {
        t.Fatal("Access token missing")
    }

    // 7. Test UserInfo
    req = httptest.NewRequest(http.MethodGet, "/userinfo", nil)
    req.Header.Set("Authorization", "Bearer "+accessToken)
    rec = httptest.NewRecorder()
    c = e.NewContext(req, rec)

    claims, _ := utils.ValidateJWT(accessToken)
    c.Set("user_id", claims.Subject)

    if err := oauthHandler.UserInfo(c); err != nil {
         t.Errorf("UserInfo failed: %v", err)
    }
    if rec.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", rec.Code)
    }
    var userResp map[string]interface{}
    json.Unmarshal(rec.Body.Bytes(), &userResp)
    if userResp["name"] != "testuser" {
        t.Errorf("Expected username testuser, got %v", userResp["name"])
    }

	// 8. Test Replay Attack
	req = httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	err = oauthHandler.Token(c)
	if err == nil {
		t.Errorf("Expected error for replay attack")
	} else {
        if httpErr, ok := err.(*echo.HTTPError); ok {
            if httpErr.Code != http.StatusBadRequest {
                 t.Errorf("Expected 400, got %d", httpErr.Code)
            }
        } else {
            t.Errorf("Expected HTTPError, got %v", err)
        }
    }
}
