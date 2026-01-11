package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"strconv"
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userService.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	token, err := h.userService.GenerateSessionToken(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate session token")
	}

	// Set session cookie
	c.SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	})

	// Check for "next" query parameter
	next := c.QueryParam("next")
	if next != "" {
		// Just returning a JSON with redirect info for now, or we could redirect if this was a form post
		// Since the request is likely JSON (Bind), we return JSON.
		// If the client expects a redirect, they should handle the "redirect_to" field.
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":     "Login successful",
			"user_id":     strconv.FormatUint(uint64(user.ID), 10),
			"redirect_to": next,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
	})
}