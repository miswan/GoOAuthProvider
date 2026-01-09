package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
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

	// Check for session
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		// Not authenticated, redirect to login
		// Include the current URL as redirect_to so user comes back here after login
		loginURL := "/login?redirect_to=" + c.Request().RequestURI
		return c.Redirect(http.StatusFound, loginURL)
	}

	// Parse user ID safely
	userIDStr, ok := userIDVal.(string)
	if !ok {
		// Should not happen if middleware is correct
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user session")
	}
	userID64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user ID in session")
	}
	userID := uint(userID64)

	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		userID,
		req.RedirectURI,
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Redirect back to client with code
	redirectURL, err := utils.BuildRedirectURL(req.RedirectURI, map[string]string{
		"code":  code,
		"state": req.State,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid redirect URI")
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

	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
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