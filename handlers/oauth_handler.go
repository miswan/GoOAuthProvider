package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"oauth2-provider/models"
	"strconv"
)

type OAuthServiceInterface interface {
	ValidateAuthorizationRequest(req *models.AuthorizationRequest) error
	GenerateAuthorizationCode(clientID string, userID uint, codeChallenge, codeChallengeMethod string) (string, error)
	ExchangeToken(req *models.TokenRequest) (string, string, error)
}

type OAuthHandler struct {
	oauthService OAuthServiceInterface
}

func NewOAuthHandler(oauthService OAuthServiceInterface) *OAuthHandler {
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

	// For simplicity, assuming user is already authenticated
	// In real implementation, check session and show login/consent page
	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		1, // Temporary userID for testing
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Construct redirect URL
	redirectURL, err := url.Parse(req.RedirectURI)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid redirect URI")
	}

	query := redirectURL.Query()
	query.Set("code", code)
	query.Set("state", req.State)
	redirectURL.RawQuery = query.Encode()

	return c.Redirect(http.StatusFound, redirectURL.String())
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