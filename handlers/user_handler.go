package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/utils"
	"strconv"
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
	contentType := c.Request().Header.Get("Content-Type")

	// Check if it's form data (browser login)
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		username := c.FormValue("username")
		password := c.FormValue("password")
		returnTo := c.FormValue("return_to")

		req := &models.UserLogin{
			Username: username,
			Password: password,
		}

		user, err := h.userService.Login(req)
		if err != nil {
			return c.HTML(http.StatusUnauthorized, "<h1>Login Failed</h1><p>Invalid credentials. <a href='/login'>Try again</a></p>")
		}

		// Generate Token
		token, err := utils.GenerateJWT(user.ID, 24*time.Hour)
		if err != nil {
			return c.HTML(http.StatusInternalServerError, "Error generating token")
		}

		// Set Cookie
		cookie := new(http.Cookie)
		cookie.Name = "auth_token"
		cookie.Value = token
		cookie.Expires = time.Now().Add(24 * time.Hour)
		cookie.HttpOnly = true
		cookie.Path = "/"
		c.SetCookie(cookie)

		if returnTo != "" {
			return c.Redirect(http.StatusFound, returnTo)
		}
		return c.Redirect(http.StatusFound, "/")
	}

	// API Login (JSON)
	req := new(models.UserLogin)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userService.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	token, _ := utils.GenerateJWT(user.ID, 24*time.Hour)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
		"token":   token,
	})
}
