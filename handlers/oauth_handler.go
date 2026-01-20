package handlers

import (
	"net/http"
	"net/url"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type OAuthHandler struct {
	oauthService *services.OAuthService
}

func NewOAuthHandler(oauthService *services.OAuthService) *OAuthHandler {
	return &OAuthHandler{oauthService: oauthService}
}

func (h *OAuthHandler) Authorize(c echo.Context) error {
	// 1. Check Session
	var userID uint
	cookie, err := c.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		// Not authenticated, redirect to login with next param
		// We use RequestURI to capture query params (client_id, etc.)
		return c.Redirect(http.StatusFound, "/login?next="+url.QueryEscape(c.Request().RequestURI))
	} else {
		token, err := utils.ValidatePaseto(cookie.Value)
		if err != nil {
			return c.Redirect(http.StatusFound, "/login?next="+url.QueryEscape(c.Request().RequestURI))
		}
		uid, _ := strconv.ParseUint(token.Subject, 10, 64)
		userID = uint(uid)
	}

	// 2. Bind Request
	req := new(models.AuthorizationRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 3. Validate Request
	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 4. Generate Code
	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		userID,
		req.RedirectURI, // Pass RedirectURI
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 5. Redirect to Client
	// Append code and state to the redirect_uri
	// Need to handle if redirect_uri already has params? Standard OAuth usually assumes base URI.
	// But let's be safe with string concatenation for now as typical usage.
	// A proper implementation would parse the redirectURI and add params.
	targetURL, _ := url.Parse(req.RedirectURI)
	q := targetURL.Query()
	q.Set("code", code)
	if req.State != "" {
		q.Set("state", req.State)
	}
	targetURL.RawQuery = q.Encode()

	return c.Redirect(http.StatusFound, targetURL.String())
}

func (h *OAuthHandler) Token(c echo.Context) error {
	req := new(models.TokenRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accessToken, refreshToken, err := h.oauthService.ExchangeToken(req)
	if err != nil {
		// OAuth 2.0 error response usually JSON with "error" field
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":             "invalid_request",
			"error_description": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
	})
}

func (h *OAuthHandler) UserInfo(c echo.Context) error {
	// UserID is set by middleware
	userID := c.Get("user_id").(string)

	// In a real app, fetch user details from UserService/DB
	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub": userID,
		// "name": "...",
	})
}
