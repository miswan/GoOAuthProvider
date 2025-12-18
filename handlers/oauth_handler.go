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

	// Check for user authentication
	cookie, err := c.Cookie("session_token")
	if err != nil {
		// In a real app, we would redirect to login page here
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	claims, err := utils.ValidateJWT(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user session")
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID in session")
	}

	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		uint(userID),
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Redirect to the client's redirect_uri with code and state
	u, err := url.Parse(req.RedirectURI)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid redirect URI")
	}
	q := u.Query()
	q.Set("code", code)
	q.Set("state", req.State)
	u.RawQuery = q.Encode()

	return c.Redirect(http.StatusFound, u.String())
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
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub":      user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}