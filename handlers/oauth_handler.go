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

	// Check Session
	cookie, err := c.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		next := "/authorize?" + c.QueryString()
		return c.Redirect(http.StatusFound, "/login?next="+url.QueryEscape(next))
	}

    // Validate Token
    claims, err := utils.ValidatePaseto(cookie.Value)
    if err != nil {
		next := "/authorize?" + c.QueryString()
        return c.Redirect(http.StatusFound, "/login?next="+url.QueryEscape(next))
    }

    userID, _ := strconv.ParseUint(claims.Subject, 10, 64)

	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		uint(userID),
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

    // Redirect to Client
    redirectURL, _ := url.Parse(req.RedirectURI)
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
        // OAuth2 error response format
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

    c.Response().Header().Set("Cache-Control", "no-store")
    c.Response().Header().Set("Pragma", "no-cache")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
	})
}

func (h *OAuthHandler) UserInfo(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
    if err != nil {
        return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID")
    }

    user, err := h.oauthService.GetUser(uint(userID))
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "User not found")
    }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub":   strconv.FormatUint(uint64(user.ID), 10),
		"name":  user.Username,
		"email": user.Email,
	})
}
