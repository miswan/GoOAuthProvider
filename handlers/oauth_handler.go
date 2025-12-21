package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"strconv"
)

type OAuthHandler struct {
	oauthService *services.OAuthService
}

func NewOAuthHandler(oauthService *services.OAuthService) *OAuthHandler {
	return &OAuthHandler{oauthService: oauthService}
}

func (h *OAuthHandler) Authorize(c echo.Context) error {
	// Check session
	cookie, err := c.Cookie("session_token")
	if err != nil {
		// Not logged in, redirect to login
		return c.Redirect(http.StatusFound, "/login?return_to="+c.Request().RequestURI)
	}

	_, err = services.ValidateSessionToken(cookie.Value)
	if err != nil {
		// Invalid session
		return c.Redirect(http.StatusFound, "/login?return_to="+c.Request().RequestURI)
	}

	req := new(models.AuthorizationRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Render consent page
	// In a real app, we would also check if consent was already given
	client := h.oauthService.GetClient(req.ClientID)
	if client == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid client_id")
	}

	return c.Render(http.StatusOK, "consent.html", map[string]interface{}{
		"ClientName":          client.ClientID, // Or a name field if it exists
		"ClientID":            req.ClientID,
		"State":               req.State,
		"RedirectURI":         req.RedirectURI,
		"CodeChallenge":       req.CodeChallenge,
		"CodeChallengeMethod": req.CodeChallengeMethod,
	})
}

func (h *OAuthHandler) ApproveConsent(c echo.Context) error {
	// Check session again
	cookie, err := c.Cookie("session_token")
	if err != nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	claims, err := services.ValidateSessionToken(cookie.Value)
	if err != nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	userID, _ := strconv.ParseUint(claims.Subject, 10, 64)

	allow := c.FormValue("allow")
	clientID := c.FormValue("client_id")
	redirectURI := c.FormValue("redirect_uri")
	state := c.FormValue("state")

	if allow != "true" {
		// Denied - Redirect with error
		redirectURL := redirectURI + "?error=access_denied"
		if state != "" {
			redirectURL += "&state=" + state
		}
		return c.Redirect(http.StatusFound, redirectURL)
	}

	codeChallenge := c.FormValue("code_challenge")
	codeChallengeMethod := c.FormValue("code_challenge_method")

	// Validate again (simplified for now)
	req := &models.AuthorizationRequest{
		ClientID:            clientID,
		RedirectURI:         redirectURI,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	code, err := h.oauthService.GenerateAuthorizationCode(
		clientID,
		uint(userID),
		codeChallenge,
		codeChallengeMethod,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Redirect to redirect_uri with code and state
	redirectURL := redirectURI + "?code=" + code
	if state != "" {
		redirectURL += "&state=" + state
	}

	return c.Redirect(http.StatusFound, redirectURL)
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