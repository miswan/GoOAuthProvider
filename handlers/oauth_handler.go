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
	// Check for auth cookie
	cookie, err := c.Cookie("auth_token")
	if err != nil {
		// Redirect to login
		return c.Redirect(http.StatusFound, "/login?continue_to="+url.QueryEscape(c.Request().RequestURI))
	}

	claims, err := utils.ValidateJWT(cookie.Value)
	if err != nil {
		return c.Redirect(http.StatusFound, "/login?continue_to="+url.QueryEscape(c.Request().RequestURI))
	}

	userID, _ := strconv.ParseUint(claims.Subject, 10, 64)

	req := new(models.AuthorizationRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.oauthService.ValidateAuthorizationRequest(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		uint(userID),
		req.CodeChallenge,
		req.CodeChallengeMethod,
		req.RedirectURI,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Redirect with code and state
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

	user := h.userService.GetUser(uint(userID))
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub":   userID,
		"name":  user.Username,
		"email": user.Email,
	})
}
