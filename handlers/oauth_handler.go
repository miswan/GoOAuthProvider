package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"oauth2-provider/models"
	"oauth2-provider/services"
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

	// Check if user is authenticated via session/cookie
	// Using standard JWT cookie "session_token" set by UserHandler
	cookie, err := c.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		// User is not authenticated, redirect to login page
		// Construct the "next" URL to return to this authorization request
		q := c.Request().URL.Query()
		nextURL := fmt.Sprintf("/authorize?%s", q.Encode())
		return c.Redirect(http.StatusFound, "/login?next="+url.QueryEscape(nextURL))
	}

	// Validate session token
	claims, err := services.ValidateSessionToken(cookie.Value)
	if err != nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	userID := claims.UserID

	// Generate authorization code
	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		req.RedirectURI,
		userID,
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Redirect back to the client application with the code and state
	redirectURL, err := url.Parse(req.RedirectURI)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid redirect URI")
	}

	q := redirectURL.Query()
	q.Set("code", code)
	if req.State != "" {
		q.Set("state", req.State)
	}
	redirectURL.RawQuery = q.Encode()

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
	// user_id is set by middleware.JWTAuth
	userID := c.Get("user_id").(uint)

	// In a real app, fetch user details from DB
	// For now, return basic info
	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub":   userID,
		"name":  fmt.Sprintf("User %d", userID),
		"email": fmt.Sprintf("user%d@example.com", userID),
	})
}
