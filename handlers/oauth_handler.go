package handlers

import (
	"html/template"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/utils"
	"strconv"
)

type OAuthHandler struct {
	oauthService *services.OAuthService
}

func NewOAuthHandler(oauthService *services.OAuthService) *OAuthHandler {
	return &OAuthHandler{oauthService: oauthService}
}

func (h *OAuthHandler) Authorize(c echo.Context) error {
	req := new(models.AuthorizationRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check for session cookie
	cookie, err := c.Cookie("session_token")
	if err != nil {
		// Redirect to login
		q := c.Request().URL.Query()
		return c.Redirect(http.StatusFound, "/login?redirect_to="+url.QueryEscape("/authorize?"+q.Encode()))
	}

	claims, err := utils.ValidateJWT(cookie.Value)
	if err != nil {
		// Invalid session, redirect to login
		q := c.Request().URL.Query()
		return c.Redirect(http.StatusFound, "/login?redirect_to="+url.QueryEscape("/authorize?"+q.Encode()))
	}

	// Show Consent Page
	_ = claims.Subject

	const tpl = `
<!DOCTYPE html>
<html>
<head>
    <title>Authorize Access</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; }
        .container { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); width: 400px; }
        h2 { text-align: center; color: #1a73e8; }
        p { text-align: center; }
        .actions { display: flex; justify-content: space-between; margin-top: 2rem; }
        button { padding: 0.5rem 1rem; border: none; border-radius: 4px; cursor: pointer; font-weight: bold; }
        .allow { background: #1a73e8; color: white; }
        .deny { background: #e0e0e0; color: #333; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Authorize Application</h2>
        <p>Client <b>{{.ClientID}}</b> is requesting access to your account.</p>
        <p>Scopes: {{.Scope}}</p>
        <form action="/approve" method="POST">
            <input type="hidden" name="client_id" value="{{.ClientID}}">
            <input type="hidden" name="response_type" value="{{.ResponseType}}">
            <input type="hidden" name="redirect_uri" value="{{.RedirectURI}}">
            <input type="hidden" name="scope" value="{{.Scope}}">
            <input type="hidden" name="state" value="{{.State}}">
            <input type="hidden" name="code_challenge" value="{{.CodeChallenge}}">
            <input type="hidden" name="code_challenge_method" value="{{.CodeChallengeMethod}}">
            <div class="actions">
                <button type="button" class="deny" onclick="window.history.back()">Deny</button>
                <button type="submit" class="allow">Allow</button>
            </div>
        </form>
    </div>
</body>
</html>`

	t, err := template.New("authorize").Parse(tpl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return t.Execute(c.Response().Writer, req)
}

func (h *OAuthHandler) Approve(c echo.Context) error {
	// Verify session again
	cookie, err := c.Cookie("session_token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Session expired")
	}
	claims, err := utils.ValidateJWT(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid session")
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID in session")
	}

	// Get form params
	clientID := c.FormValue("client_id")
	redirectURI := c.FormValue("redirect_uri")
	scope := c.FormValue("scope")
	state := c.FormValue("state")
	codeChallenge := c.FormValue("code_challenge")
	codeChallengeMethod := c.FormValue("code_challenge_method")
	responseType := c.FormValue("response_type")

	// Validate Request again to ensure redirect_uri is still valid for client
	// This prevents open redirect attacks if the user modified the hidden field
	req := &models.AuthorizationRequest{
		ClientID:            clientID,
		RedirectURI:         redirectURI,
		ResponseType:        responseType,
		Scope:               scope,
		State:               state,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Generate Code
	code, err := h.oauthService.GenerateAuthorizationCode(
		clientID,
		uint(userID),
		codeChallenge,
		codeChallengeMethod,
		scope,
		redirectURI,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Redirect to callback
	callbackURL, err := url.Parse(redirectURI)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid redirect URI")
	}
	q := callbackURL.Query()
	q.Set("code", code)
	q.Set("state", state)
	callbackURL.RawQuery = q.Encode()

	return c.Redirect(http.StatusFound, callbackURL.String())
}

func (h *OAuthHandler) Token(c echo.Context) error {
	req := new(models.TokenRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accessToken, refreshToken, err := h.oauthService.ExchangeToken(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    "3600",
	})
}

func (h *OAuthHandler) UserInfo(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub":   userID,
		"name":  "John Doe",
		"email": "john@example.com",
	})
}
