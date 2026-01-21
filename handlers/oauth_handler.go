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

	// Validate Request
	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check Authentication via Cookie
	cookie, err := c.Cookie("auth_token")
	if err != nil || cookie.Value == "" {
		// Redirect to login
		// Construct return URL
		q := c.Request().URL.Query()
		returnURL := "/authorize?" + q.Encode()
		return c.Redirect(http.StatusFound, "/login?redirect_to="+url.QueryEscape(returnURL))
	}

	// Validate JWT from cookie
	claims, err := utils.ValidateJWT(cookie.Value)
	if err != nil {
		return c.Redirect(http.StatusFound, "/login?redirect_to="+url.QueryEscape(c.Request().URL.String()))
	}

	userID, _ := strconv.ParseUint(claims.Subject, 10, 64)

	// Generate Code
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

	// Construct Redirect URL with Code and State
	redirectURL, err := url.Parse(req.RedirectURI)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid redirect uri")
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

	// mock user info for now, or fetch from service
	// In a real app we would use userService to fetch user details

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub":   userIDStr,
		"name":  "User " + userIDStr,
		// "email": "...", // need to fetch user
	})
}
