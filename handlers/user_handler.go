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

	// If Form Post, redirect to login
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		return c.Redirect(http.StatusFound, "/login")
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

	token, err := utils.GeneratePaseto(user.ID, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// Handle Redirect
	next := c.FormValue("next")
	if next == "" {
		next = c.QueryParam("next")
	}

	contentType := c.Request().Header.Get("Content-Type")
	// If JSON, return JSON
	if contentType == "application/json" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Login successful",
			"user_id": strconv.FormatUint(uint64(user.ID), 10),
			"token":   token,
		})
	}

	// Default to Redirect for Browser/Form
	if next != "" && utils.IsValidRedirect(next) {
		return c.Redirect(http.StatusFound, next)
	}

	return c.Redirect(http.StatusFound, "/")
}
