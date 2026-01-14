package handlers

import (
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
	req := new(models.AuthorizationRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check for session token
	cookie, err := c.Cookie("session_token")
	if err != nil {
		// Redirect to login with next param
		loginURL := "/login?next=" + url.QueryEscape(c.Request().URL.String())
		return c.Redirect(http.StatusFound, loginURL)
	}

	// Validate session (simple validation for now, relying on utils)
	// In a real app, you might want to look up the user in the DB
	claims, err := services.ValidateSessionToken(cookie.Value)
	if err != nil {
		return c.Redirect(http.StatusFound, "/login?next="+url.QueryEscape(c.Request().URL.String()))
	}

	userID, _ := strconv.ParseUint(claims.Subject, 10, 64)

	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		uint(userID),
		req.RedirectURI,
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Construct redirect URL safely
	u, err := url.Parse(req.RedirectURI)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid redirect URI")
	}
	q := u.Query()
	q.Set("code", code)
	q.Set("state", req.State)
	u.RawQuery = q.Encode()

	return c.Redirect(http.StatusFound, u.String())
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