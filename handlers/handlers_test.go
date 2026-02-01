package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

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

	// Create Client
	client := &models.Client{
		ClientID:     "test-client",
		Secret:       "test-secret",
		RedirectURIs: []string{"http://localhost:8080/callback"},
		GrantTypes:   []string{"authorization_code"},
	}
	// Need to initialize ID for memory storage if not handled
	// models.Client is gorm.Model, ID is uint.
	client.ID = 1
	store.StoreClient(client)

	// Create User
	pwd, _ := utils.HashPassword("password123")
	user := &models.User{
		Username: "testuser",
		Password: pwd,
		Email:    "test@example.com",
	}
	user.ID = 1
	store.StoreUser(user)

	// Init Services & Handlers
	oauthService := services.NewOAuthService(store)
	userService := services.NewUserService(store)

	userHandler := handlers.NewUserHandler(userService)
	oauthHandler := handlers.NewOAuthHandler(oauthService)
	// htmlHandler := handlers.NewHTMLHandler() // Not tested here directly

	// 1. Test Login (POST)
	f := make(url.Values)
	f.Set("username", "testuser")
	f.Set("password", "password123")
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := userHandler.Login(c); err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Login status: %d", rec.Code)
	}

	// Extract Cookie
	cookies := rec.Result().Cookies()
	var authToken string
	for _, ck := range cookies {
		if ck.Name == "auth_token" {
			authToken = ck.Value
		}
	}
	if authToken == "" {
		t.Fatal("Auth token cookie not set")
	}

	// 2. Test Authorize (GET)
	q := make(url.Values)
	q.Set("client_id", "test-client")
	q.Set("redirect_uri", "http://localhost:8080/callback")
	q.Set("response_type", "code")
	q.Set("state", "xyz")
	q.Set("code_challenge", "secret")
	q.Set("code_challenge_method", "plain")

	reqAuth := httptest.NewRequest(http.MethodGet, "/authorize?"+q.Encode(), nil)
	reqAuth.Header.Set("Cookie", "auth_token="+authToken)
	recAuth := httptest.NewRecorder()
	cAuth := e.NewContext(reqAuth, recAuth)

	if err := oauthHandler.Authorize(cAuth); err != nil {
		t.Fatalf("Authorize failed: %v", err)
	}

	if recAuth.Code != http.StatusFound {
		t.Fatalf("Authorize status: %d, body: %s", recAuth.Code, recAuth.Body.String())
	}

	loc := recAuth.Header().Get("Location")
	if !strings.Contains(loc, "code=") {
		t.Fatalf("Redirect location missing code: %s", loc)
	}

	// Parse code from Location
	locURL, _ := url.Parse(loc)
	code := locURL.Query().Get("code")

	// 3. Test Token Exchange
	fToken := make(url.Values)
	fToken.Set("grant_type", "authorization_code")
	fToken.Set("code", code)
	fToken.Set("client_id", "test-client")
	fToken.Set("redirect_uri", "http://localhost:8080/callback")
	fToken.Set("code_verifier", "secret")

	reqToken := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(fToken.Encode()))
	reqToken.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	recToken := httptest.NewRecorder()
	cToken := e.NewContext(reqToken, recToken)

	if err := oauthHandler.Token(cToken); err != nil {
		t.Fatalf("Token failed: %v", err)
	}

	if recToken.Code != http.StatusOK {
		t.Fatalf("Token status: %d, body: %s", recToken.Code, recToken.Body.String())
	}

	if !strings.Contains(recToken.Body.String(), "access_token") {
		t.Fatal("Response missing access_token")
	}
}
