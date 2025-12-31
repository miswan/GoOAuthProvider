package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"strconv"
	"strings"
)

type OAuthHandler struct {
	oauthService *services.OAuthService
}

func NewOAuthHandler(oauthService *services.OAuthService) *OAuthHandler {
	return &OAuthHandler{oauthService: oauthService}
}

func (h *OAuthHandler) Authorize(c echo.Context) error {
	// Check if user is authenticated via session
	// We expect middleware to set user_id if valid session exists
	userID := c.Get("user_id")
	if userID == nil {
		// Not authenticated, redirect to login with return_to
		returnTo := c.Request().RequestURI // Contains path and query params
		return c.Redirect(http.StatusFound, "/login?redirect_to="+returnTo)
	}

	// Convert userID to uint
	var uid uint
	if idStr, ok := userID.(string); ok {
		// Assuming ID is stored as string in JWT claims
		idUint, _ := strconv.ParseUint(idStr, 10, 64)
		uid = uint(idUint)
	} else if idUint, ok := userID.(uint); ok {
		uid = idUint
	} else {
		// Fallback for numeric ID if stored that way in context
		// This handles cases where context might have different type
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user context")
	}

	req := new(models.AuthorizationRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// In a full implementation, we would show a consent screen here.
	// For this task, we auto-approve if logged in.

	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		uid,
		req.RedirectURI,
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Redirect back to client with code and state
	// Ensure we handle existing query parameters correctly
	separator := "?"
	if strings.Contains(req.RedirectURI, "?") {
		separator = "&"
	}

	redirectURL := req.RedirectURI + separator + "code=" + code
	if req.State != "" {
		redirectURL += "&state=" + req.State
	}

	return c.Redirect(http.StatusFound, redirectURL)
}

func (h *OAuthHandler) Token(c echo.Context) error {
	req := new(models.TokenRequest)
	// Bind supports JSON and Form (application/x-www-form-urlencoded) due to updated struct tags
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