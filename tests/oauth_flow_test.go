package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"oauth2-provider/handlers"
	"oauth2-provider/middleware"
	"oauth2-provider/services"
	"oauth2-provider/storage"
)

func TestOAuthFlow(t *testing.T) {
	// Setup
	store := storage.NewMemoryStorage()
	oauthService := services.NewOAuthService(store)
	userService := services.NewUserService(store)
	clientService := services.NewClientService(store)

	oauthHandler := handlers.NewOAuthHandler(oauthService, userService)
	userHandler := handlers.NewUserHandler(userService)
	clientHandler := handlers.NewClientHandler(clientService)

	e := echo.New()

	// 1. Register User
	t.Run("Register User", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"testuser","password":"password123","email":"test@example.com"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, userHandler.Register(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	// 2. Login User
	var authToken string
	t.Run("Login User", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"testuser","password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, userHandler.Login(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			// Check cookie
			for _, cookie := range rec.Result().Cookies() {
				if cookie.Name == "auth_token" {
					authToken = cookie.Value
				}
			}
			assert.NotEmpty(t, authToken)
		}
	})

	// 3. Register Client
	var clientID, clientSecret string
	t.Run("Register Client", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/client/register", strings.NewReader(`{"redirect_uris":["http://localhost:8080/callback"]}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, clientHandler.Register(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			var resp map[string]string
			json.Unmarshal(rec.Body.Bytes(), &resp)
			clientID = resp["client_id"]
			clientSecret = resp["client_secret"]
			assert.NotEmpty(t, clientID)
			assert.NotEmpty(t, clientSecret)
		}
	})

	// 4. Authorize (Unauthenticated)
	t.Run("Authorize Unauthenticated", func(t *testing.T) {
		q := make(url.Values)
		q.Set("client_id", clientID)
		q.Set("redirect_uri", "http://localhost:8080/callback")
		q.Set("response_type", "code")
		q.Set("code_challenge", "challenge")
		q.Set("code_challenge_method", "plain")

		req := httptest.NewRequest(http.MethodGet, "/authorize?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// It returns error because Redirect writes header and status code, Echo might wrap it.
		// Actually c.Redirect returns error.
		err := oauthHandler.Authorize(c)
		// Check redirect location
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Contains(t, rec.Header().Get("Location"), "/login")
		assert.NoError(t, err)
	})

	// 5. Authorize (Authenticated)
	var authCode string
	t.Run("Authorize Authenticated", func(t *testing.T) {
		q := make(url.Values)
		q.Set("client_id", clientID)
		q.Set("redirect_uri", "http://localhost:8080/callback")
		q.Set("response_type", "code")
		q.Set("code_challenge", "challenge")
		q.Set("code_challenge_method", "plain")
		q.Set("state", "mystate")

		req := httptest.NewRequest(http.MethodGet, "/authorize?"+q.Encode(), nil)
		req.AddCookie(&http.Cookie{Name: "auth_token", Value: authToken})
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, oauthHandler.Authorize(c)) {
			assert.Equal(t, http.StatusFound, rec.Code)
			loc := rec.Header().Get("Location")
			assert.Contains(t, loc, "http://localhost:8080/callback")

			u, _ := url.Parse(loc)
			authCode = u.Query().Get("code")
			assert.NotEmpty(t, authCode)
			assert.Equal(t, "mystate", u.Query().Get("state"))
		}
	})

	// 6. Exchange Token
	var accessToken string
	t.Run("Exchange Token", func(t *testing.T) {
		// Needs JSON body for c.Bind or Form? TokenRequest has JSON tags.
		// Note: The memory says: "Go structs used for binding HTTP form data ... must include form tags"
		// models.TokenRequest might need form tags if we send form data.
		// But let's send JSON for now as handlers support Bind (which supports JSON).
		// Wait, OAuth specs usually use Form POST.
		// Let's check TokenRequest struct.

		body := map[string]string{
			"grant_type":    "authorization_code",
			"code":          authCode,
			"redirect_uri":  "http://localhost:8080/callback", // Must match
			"client_id":     clientID,
			"client_secret": clientSecret,
			"code_verifier": "challenge", // plain
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(string(bodyBytes)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, oauthHandler.Token(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var resp map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &resp)
			accessToken = resp["access_token"].(string)
			assert.NotEmpty(t, accessToken)
		}
	})

	// 7. User Info
	t.Run("User Info", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Setup middleware chain
		h := middleware.JWTAuth(oauthHandler.UserInfo)

		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var resp map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.Equal(t, "testuser", resp["name"])
			assert.Equal(t, "test@example.com", resp["email"])
		}
	})
}
