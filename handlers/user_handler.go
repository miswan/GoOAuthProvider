package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/utils"
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
		next := c.QueryParam("next")
		return c.Render(http.StatusBadRequest, "register.html", map[string]interface{}{
			"Error": "Invalid request: " + err.Error(),
			"Next":  next,
		})
	}

	if err := h.userService.Register(req); err != nil {
		next := c.QueryParam("next")
		return c.Render(http.StatusBadRequest, "register.html", map[string]interface{}{
			"Error": err.Error(),
			"Next":  next,
		})
	}

	// Check if this is an API request or Form submission
	if c.Request().Header.Get("Content-Type") == "application/json" {
		return c.JSON(http.StatusCreated, map[string]string{
			"message": "User registered successfully",
		})
	}

	// For form submission, redirect to login
	redirectURL := "/login"
	next := c.QueryParam("next")
	if next != "" {
		redirectURL += "?next=" + next
	}
	return c.Redirect(http.StatusSeeOther, redirectURL)
}

func (h *UserHandler) Login(c echo.Context) error {
	req := new(models.UserLogin)
	if err := c.Bind(req); err != nil {
		next := c.QueryParam("next")
		return c.Render(http.StatusBadRequest, "login.html", map[string]interface{}{
			"Error": "Invalid request: " + err.Error(),
			"Next":  next,
		})
	}

	user, err := h.userService.Login(req)
	if err != nil {
		next := c.QueryParam("next")
		return c.Render(http.StatusUnauthorized, "login.html", map[string]interface{}{
			"Error": "Invalid credentials",
			"Next":  next,
		})
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	// Set cookie
	cookie := new(http.Cookie)
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	// Check if this is an API request or Form submission
	if c.Request().Header.Get("Content-Type") == "application/json" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Login successful",
			"token":   token,
		})
	}

	// For form submission, redirect
	next := c.QueryParam("next")
	if next != "" {
		return c.Redirect(http.StatusSeeOther, next)
	}
	return c.Redirect(http.StatusSeeOther, "/")
}
