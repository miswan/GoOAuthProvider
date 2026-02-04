package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/utils"
	"strconv"
	"time"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c echo.Context) error {
	req := new(models.UserRegister)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.userService.Register(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	req := new(models.UserLogin)
	if err := c.Bind(req); err != nil {
		// Try to parse form values if JSON binding fails (handled by Bind usually but fallback just in case)
		// Actually echo.Bind handles Content-Type check.
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userService.Login(req)
	if err != nil {
		// Check for return_to parameter to handle login failure redirection
		returnTo := c.FormValue("return_to")
		if returnTo != "" && utils.IsValidReturnTo(returnTo) {
			return c.Redirect(http.StatusFound, "/login?error=Invalid credentials&return_to="+returnTo)
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Set auth_token cookie
	c.SetCookie(&http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
	})

	// Handle return_to redirect
	returnTo := c.FormValue("return_to")
	if returnTo != "" && utils.IsValidReturnTo(returnTo) {
		return c.Redirect(http.StatusFound, returnTo)
	}

	// Default JSON response for API usage
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
		"token":   token,
	})
}
