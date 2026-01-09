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
	if c.Request().Method == http.MethodGet {
		redirectTo := c.QueryParam("redirect_to")
		html := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Login</title>
			</head>
			<body>
				<h2>Login</h2>
				<form action="/login" method="POST">
					<input type="hidden" name="redirect_to" value="` + redirectTo + `">
					<div>
						<label>Username:</label>
						<input type="text" name="username" required>
					</div>
					<div>
						<label>Password:</label>
						<input type="password" name="password" required>
					</div>
					<button type="submit">Login</button>
				</form>
				<p><a href="/register">Register</a></p>
			</body>
			</html>
		`
		return c.HTML(http.StatusOK, html)
	}

	req := new(models.UserLogin)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userService.Login(req)
	if err != nil {
		return c.HTML(http.StatusUnauthorized, "Invalid credentials")
	}

	// Generate session token (JWT)
	// Using a longer expiration for session cookies (e.g., 24 hours)
	token, err := h.userService.GenerateSessionToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate session")
	}

	// Set session cookie
	cookie := new(http.Cookie)
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	// cookie.Secure = true // Enable in production with HTTPS
	c.SetCookie(cookie)

	redirectTo := c.FormValue("redirect_to")
	if redirectTo != "" {
		return c.Redirect(http.StatusFound, redirectTo)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
	})
}