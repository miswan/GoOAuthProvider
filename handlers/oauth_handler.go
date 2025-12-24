package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
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
	// Check if user is authenticated
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		// Redirect to login with return_to
		// We should construct the full URL including query params
		req := c.Request()

		// c.Request().RequestURI gives the full path + query.
		returnTo := req.RequestURI
		return c.Redirect(http.StatusFound, fmt.Sprintf("/login?return_to=%s", url.QueryEscape(returnTo)))
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user ID")
	}

	req := new(models.AuthorizationRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Generate authorization code
	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		uint(userID),
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Redirect to redirect_uri with code and state
	// Parse the redirect URI to handle existing query parameters
	parsedURI, err := url.Parse(req.RedirectURI)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid redirect URI")
	}

	query := parsedURI.Query()
	query.Set("code", code)
	if req.State != "" {
		query.Set("state", req.State)
	}
	parsedURI.RawQuery = query.Encode()

	return c.Redirect(http.StatusFound, parsedURI.String())
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