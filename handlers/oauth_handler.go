package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"strconv"
	"strings"
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
		return c.HTML(http.StatusBadRequest, "Invalid request: "+err.Error())
	}

	// Validate request parameters early
	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return c.HTML(http.StatusBadRequest, "Authorization Error: "+err.Error())
	}

	// Check authentication
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		// Not authenticated, redirect to login
		returnTo := c.Request().URL.Path
		if c.QueryString() != "" {
			returnTo += "?" + c.QueryString()
		}
		return c.Redirect(http.StatusFound, "/login?return_to="+url.QueryEscape(returnTo))
	}

	// Parse user ID
	userIDStr := userIDVal.(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	// Generate authorization code
	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		uint(userID),
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, "Server Error: "+err.Error())
	}

	// Construct redirect URL
	redirectURL := req.RedirectURI
	separator := "?"
	if strings.Contains(redirectURL, "?") {
		separator = "&"
	}
	redirectURL += separator + "code=" + code
	if req.State != "" {
		redirectURL += "&state=" + req.State
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
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub":   user.ID,
		"name":  user.Username,
		"email": user.Email,
	})
}
