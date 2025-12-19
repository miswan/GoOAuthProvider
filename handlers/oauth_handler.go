package handlers

import (
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

	// Check session
	cookie, err := c.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		// Redirect to login
		// We need to construct the full return URL.
		// Since we are in the handler, `c.Request().RequestURI` should contain the path and query string.
		return c.Redirect(http.StatusFound, "/login?return_to="+url.QueryEscape(c.Request().RequestURI))
	}

	claims, err := utils.ValidateJWT(cookie.Value)
	if err != nil {
		// Invalid token, redirect to login
		return c.Redirect(http.StatusFound, "/login?return_to="+url.QueryEscape(c.Request().RequestURI))
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		// If subject is not a valid user ID, treat as invalid session
		return c.Redirect(http.StatusFound, "/login?return_to="+url.QueryEscape(c.Request().RequestURI))
	}

	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		uint(userID),
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Redirect back to client with code and state
	redirectURL, err := url.Parse(req.RedirectURI)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid redirect URI")
	}

	query := redirectURL.Query()
	query.Set("code", code)
	if req.State != "" {
		query.Set("state", req.State)
	}
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
	userIDVal := c.Get("user_id")
	userIDStr, ok := userIDVal.(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user session")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user ID format")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub": userID,
	})
}
