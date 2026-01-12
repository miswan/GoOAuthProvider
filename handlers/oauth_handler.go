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

	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check for session cookie
	cookie, err := c.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		// Redirect to login
		loginURL := "/login?next=" + url.QueryEscape(c.Request().RequestURI)
		return c.Redirect(http.StatusFound, loginURL)
	}

	// Validate token
	claims, err := utils.ValidateJWT(cookie.Value)
	if err != nil {
		loginURL := "/login?next=" + url.QueryEscape(c.Request().RequestURI)
		return c.Redirect(http.StatusFound, loginURL)
	}

	userIDStr := claims.Subject
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		req.RedirectURI,
		uint(userID),
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	redirectURL, err := utils.BuildRedirectURL(req.RedirectURI, map[string]string{
		"code":  code,
		"state": req.State,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

	user, err := h.userService.GetUserByID(uint(userID))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub":   userID,
		"name":  user.Username,
		"email": user.Email,
	})
}
