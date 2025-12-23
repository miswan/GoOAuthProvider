package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
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

// LoginPage serves the HTML login form
func (h *UserHandler) LoginPage(c echo.Context) error {
	return_to := c.QueryParam("return_to")
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Login</title>
			<style>
				body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; }
				.login-container { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); width: 300px; }
				input { width: 100%%; padding: 10px; margin: 10px 0; border: 1px solid #ccc; border-radius: 4px; box-sizing: border-box; }
				button { width: 100%%; padding: 10px; background-color: #1877f2; color: white; border: none; border-radius: 4px; cursor: pointer; }
				button:hover { background-color: #166fe5; }
			</style>
		</head>
		<body>
			<div class="login-container">
				<h2>Login</h2>
				<form action="/login" method="POST">
					<input type="hidden" name="return_to" value="%s">
					<input type="text" name="username" placeholder="Username" required>
					<input type="password" name="password" placeholder="Password" required>
					<button type="submit">Log In</button>
				</form>
				<p><a href="/register">Register</a></p>
			</div>
		</body>
		</html>
	`, return_to)
	return c.HTML(http.StatusOK, html)
}

func (h *UserHandler) Login(c echo.Context) error {
	// Parse form params first to handle both JSON and Form submissions
	username := c.FormValue("username")
	password := c.FormValue("password")
	returnTo := c.FormValue("return_to")

	// If form values are empty, try binding from JSON (for API clients)
	if username == "" || password == "" {
		req := new(models.UserLogin)
		if err := c.Bind(req); err == nil {
			username = req.Username
			password = req.Password
		}
	}

	if username == "" || password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username and password are required")
	}

	user, err := h.userService.Login(&models.UserLogin{
		Username: username,
		Password: password,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// Generate session token (JWT)
	token, err := utils.GenerateJWT(user.ID, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	// Set cookie
	cookie := new(http.Cookie)
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true
	// cookie.Secure = true // Enable in production
	c.SetCookie(cookie)

	// Handle return_to redirection
	if returnTo != "" {
		// Basic security check: ensure it's a relative path to avoid open redirects
		u, err := url.Parse(returnTo)
		if err == nil && !u.IsAbs() {
			return c.Redirect(http.StatusFound, returnTo)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
		"token":   token, // Return token in body as well for API clients
	})
}
