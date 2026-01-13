package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/utils"
	"strings"
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
	// Support both JSON and Form binding
	req := new(models.UserLogin)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userService.Login(req)
	if err != nil {
		// If it's a form request (HTML login), redirect back to login with error
		if c.FormValue("username") != "" || c.FormValue("password") != "" {
			next := c.QueryParam("next")
			redirectURL := "/login?error=Invalid+credentials"
			if next != "" {
				redirectURL += "&next=" + url.QueryEscape(next)
			}
			return c.Redirect(http.StatusSeeOther, redirectURL)
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// Generate JWT (Session Token)
	// We'll use the same JWT util, assuming it puts UserID in Subject
	token, err := utils.GenerateJWT(user.ID, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	// Set Cookie
	cookie := new(http.Cookie)
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true // Important for security
	c.SetCookie(cookie)

	// Check for 'next' parameter for redirection
	next := c.QueryParam("next")
	if next != "" {
		// Validate redirect to prevent Open Redirect
		if strings.HasPrefix(next, "/") && !strings.HasPrefix(next, "//") {
			return c.Redirect(http.StatusFound, next)
		}
	}

	// Default JSON response for API clients or if no next param
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	})
}
