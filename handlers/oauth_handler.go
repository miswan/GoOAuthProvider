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
	userService  *services.UserService
}

func NewOAuthHandler(oauthService *services.OAuthService, userService *services.UserService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		userService:  userService,
	}
}

func (h *OAuthHandler) Authorize(c echo.Context) error {
	req := new(models.AuthorizationRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check authentication
	var userID uint
	cookie, err := c.Cookie("auth_token")
	if err == nil {
		claims, err := utils.ValidateJWT(cookie.Value)
		if err == nil {
			// Subject is string, parse it
			// claims.Subject should be userID string
			if uid, err := strconv.ParseUint(claims.Subject, 10, 64); err == nil {
				userID = uint(uid)
			}
		}
	}

	if userID == 0 {
		// Redirect to login
		// We should construct return_to properly
		// c.Request().RequestURI gives the relative path + query
		returnTo := c.Request().RequestURI
		if returnTo == "" {
             // If RequestURI is empty (sometimes happens in tests or depending on server/proxy), fallback
             returnTo = c.Path() + "?" + c.QueryString()
        }

		loginURL := fmt.Sprintf("/login?return_to=%s", url.QueryEscape(returnTo))
		return c.Redirect(http.StatusFound, loginURL)
	}

	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

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

	// Redirect to callback
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
	userIDStr := c.Get("user_id").(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	user, err := h.userService.GetUser(uint(userID))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub":   userID,
		"name":  user.Username,
		"email": user.Email,
	})
}
