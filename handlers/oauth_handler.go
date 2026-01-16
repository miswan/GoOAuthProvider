package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/storage"
	"oauth2-provider/utils"
	"strconv"
)

type OAuthHandler struct {
	oauthService *services.OAuthService
	store        storage.Storage
}

func NewOAuthHandler(oauthService *services.OAuthService, store storage.Storage) *OAuthHandler {
	return &OAuthHandler{oauthService: oauthService, store: store}
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
		// Not logged in, redirect to login
		nextURL := url.QueryEscape(c.Request().RequestURI)
		return c.Redirect(http.StatusFound, "/login?next="+nextURL)
	}

	// Validate session token
	claims, err := utils.ValidatePaseto(cookie.Value)
	if err != nil {
		// Invalid session, redirect to login
		nextURL := url.QueryEscape(c.Request().RequestURI)
		return c.Redirect(http.StatusFound, "/login?next="+nextURL)
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user ID in session")
	}

	// Generate authorization code
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

	// Redirect back to client with code and state
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

	user := h.store.GetUserByID(uint(userID))
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub":   user.ID,
		"name":  user.Username,
		"email": user.Email,
	})
}
