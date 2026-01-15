package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"oauth2-provider/middleware"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestOAuthFlow(t *testing.T) {
	// Setup
	store := storage.NewMemoryStorage()
	oauthService := services.NewOAuthService(store)
	userService := services.NewUserService(store)
	clientService := services.NewClientService(store)

	oauthHandler := NewOAuthHandler(oauthService)
	userHandler := NewUserHandler(userService)

	e := echo.New()

	// 1. Register User
	reqRegister := &models.UserRegister{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}
	err := userService.Register(reqRegister)
	assert.NoError(t, err)

	// 2. Register Client
	clientReq := &models.ClientRegistration{
		RedirectURIs: []string{"http://localhost:8080/callback"},
	}
	client, err := clientService.RegisterClient(clientReq)
	assert.NoError(t, err)

	// 3. Login to get Session Cookie
	loginForm := url.Values{}
	loginForm.Set("username", "testuser")
	loginForm.Set("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginForm.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = userHandler.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	assert.NotEmpty(t, rec.Result().Cookies())
	cookie := rec.Result().Cookies()[0]
	assert.Equal(t, "session_token", cookie.Name)
	sessionToken := cookie.Value

	// 4. Authorize Request
	q := make(url.Values)
	q.Set("client_id", client.ClientID)
	q.Set("redirect_uri", "http://localhost:8080/callback")
	q.Set("response_type", "code")
	q.Set("code_challenge", "plain_challenge")
	q.Set("code_challenge_method", "plain")
	q.Set("state", "xyz")

	req = httptest.NewRequest(http.MethodGet, "/authorize?"+q.Encode(), nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: sessionToken})
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	err = oauthHandler.Authorize(c)
	// c.Redirect returns nil on success usually
	if err != nil {
		assert.Fail(t, "Authorize returned error: "+err.Error())
	}
	assert.Equal(t, http.StatusFound, rec.Code)

	location := rec.Header().Get("Location")
	parsedLoc, _ := url.Parse(location)
	code := parsedLoc.Query().Get("code")
	assert.NotEmpty(t, code)

	// 5. Exchange Token
	tokenForm := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  "http://localhost:8080/callback",
		"client_id":     client.ClientID,
		"code_verifier": "plain_challenge",
	}
	jsonBody, _ := json.Marshal(tokenForm)

	req = httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(string(jsonBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	err = oauthHandler.Token(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var tokenResp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &tokenResp)
	accessToken, ok := tokenResp["access_token"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, accessToken)

	// 6. UserInfo
	req = httptest.NewRequest(http.MethodGet, "/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	// Manually wrap handler with middleware
	handler := middleware.PasetoAuth(oauthHandler.UserInfo)

	err = handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var userResp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &userResp)
	assert.Equal(t, "testuser", userResp["name"])
}
