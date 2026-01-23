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
	return &OAuthHandler{oauthService: oauthService, userService: userService}
}

func (h *OAuthHandler) Authorize(c echo.Context) error {
	req := new(models.AuthorizationRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check authentication
	var userID uint
	cookie, err := c.Cookie("auth_token")
	if err != nil {
		// Redirect to login
		q := c.Request().URL.Query().Encode()
		return c.Redirect(http.StatusFound, "/login?continue_to=/authorize?"+url.QueryEscape(q))
	}

	claims, err := utils.ValidateJWT(cookie.Value)
	if err != nil {
		// Redirect to login
		q := c.Request().URL.Query().Encode()
		return c.Redirect(http.StatusFound, "/login?continue_to=/authorize?"+url.QueryEscape(q))
	}

	uid, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user ID in token")
	}
	userID = uint(uid)

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
	userID := c.Get("user_id").(uint)

	user, err := h.userService.GetUser(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub":   fmt.Sprintf("%d", user.ID),
		"name":  user.Username,
		"email": user.Email,
	})
}