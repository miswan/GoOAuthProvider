package handlers

import (
	"fmt"
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
	// 1. Check for session cookie
	cookie, err := c.Cookie("session_token")
	var userID uint

	if err == nil {
		claims, err := utils.ValidateJWT(cookie.Value)
		if err == nil {
			uid, _ := strconv.ParseUint(claims.Subject, 10, 64)
			userID = uint(uid)
		}
	}

	// 2. If not authenticated, redirect to login with return_to
	if userID == 0 {
		currentURL := c.Request().URL.String()
		loginURL := fmt.Sprintf("/login?return_to=%s", url.QueryEscape(currentURL))
		return c.Redirect(http.StatusFound, loginURL)
	}

	req := new(models.AuthorizationRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 3. Generate auth code
	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		userID,
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 4. Redirect to callback URL
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
	userIDStr := c.Get("user_id").(string)
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID")
	}

	user, err := h.oauthService.GetUserByID(uint(userID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user")
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub":   strconv.Itoa(int(user.ID)),
		"name":  user.Username,
		"email": user.Email,
	})
}
