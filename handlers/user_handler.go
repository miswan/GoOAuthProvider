package handlers

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"strconv"
	"strings"
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
	// Handle GET request to show login form
	if c.Request().Method == http.MethodGet {
		redirectTo := c.QueryParam("redirect_to")
		// Basic XSS prevention: strict validation of redirect_to
		// It must be a relative path starting with /
		if redirectTo != "" && (!strings.HasPrefix(redirectTo, "/") || strings.HasPrefix(redirectTo, "//")) {
			redirectTo = ""
		}

		// Use html/template to prevent XSS
		tmplStr := `
			<!DOCTYPE html>
			<html>
			<head><title>Login</title></head>
			<body>
				<h2>Login</h2>
				<form method="POST" action="/login">
					<input type="hidden" name="redirect_to" value="{{.RedirectTo}}">
					<label>Username: <input type="text" name="username" required></label><br>
					<label>Password: <input type="password" name="password" required></label><br>
					<button type="submit">Login</button>
				</form>
				<p><a href="/register">Register</a></p>
			</body>
			</html>
		`
		tmpl, err := template.New("login").Parse(tmplStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Template error")
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return tmpl.Execute(c.Response().Writer, map[string]string{
			"RedirectTo": redirectTo,
		})
	}

	// Handle POST request
	req := new(models.UserLogin)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userService.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// Generate session token
	token, err := h.userService.GenerateSessionToken(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate session")
	}

	// Set cookie
	c.SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})

	// Check if there is a redirect_to param in form or query
	redirectTo := c.FormValue("redirect_to")
	if redirectTo == "" {
		redirectTo = c.QueryParam("redirect_to")
	}

	if redirectTo != "" {
		// Validate redirect_to to prevent open redirect
		if strings.HasPrefix(redirectTo, "/") && !strings.HasPrefix(redirectTo, "//") {
			return c.Redirect(http.StatusFound, redirectTo)
		}
		// If invalid redirect, just return JSON success
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
	})
}