package handlers

import (
	"html/template"
	"net/http"
	"net/url"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/utils"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
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

func (h *UserHandler) ShowLoginPage(c echo.Context) error {
	returnTo := c.QueryParam("return_to")
	action := "/login"
	if returnTo != "" {
		action += "?return_to=" + url.QueryEscape(returnTo)
	}

	errorMsg := c.QueryParam("error")

	// Use html/template to prevent XSS
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
    <style>
        body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; }
        .login-container { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); width: 300px; }
        input[type="text"], input[type="password"] { width: 100%; padding: 10px; margin: 10px 0; border: 1px solid #ddd; box-sizing: border-box; }
        input[type="submit"] { width: 100%; padding: 10px; background-color: #1877f2; color: white; border: none; border-radius: 4px; cursor: pointer; }
        input[type="submit"]:hover { background-color: #166fe5; }
        .error { color: red; font-size: 0.9em; margin-bottom: 10px; }
    </style>
</head>
<body>
<div class="login-container">
    <h2>Login</h2>
    {{if .Error}}
    <div class="error">{{.Error}}</div>
    {{end}}
    <form action="{{.Action}}" method="POST">
        <label for="username">Username</label>
        <input type="text" id="username" name="username" required>

        <label for="password">Password</label>
        <input type="password" id="password" name="password" required>

        <input type="submit" value="Log In">
    </form>
</div>
</body>
</html>
`
	t, err := template.New("login").Parse(tmpl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Template error")
	}

	data := struct {
		Action string
		Error  string
	}{
		Action: action,
		Error:  errorMsg,
	}

	return t.Execute(c.Response().Writer, data)
}

func (h *UserHandler) Login(c echo.Context) error {
	req := new(models.UserLogin)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userService.Login(req)
	if err != nil {
		// If return_to is present, imply browser flow, redirect back to login with error
		returnTo := c.QueryParam("return_to")
		if returnTo != "" {
			return c.Redirect(http.StatusFound, "/login?return_to="+url.QueryEscape(returnTo)+"&error=Invalid credentials")
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// Generate Session Token
	token, err := utils.GenerateJWT(user.ID, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate session token")
	}

	// Set Cookie
	cookie := new(http.Cookie)
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	// Check return_to and validate it to prevent Open Redirect
	returnTo := c.QueryParam("return_to")
	if returnTo != "" {
		// Ensure it's a relative path and doesn't start with // (protocol relative)
		if strings.HasPrefix(returnTo, "/") && !strings.HasPrefix(returnTo, "//") {
			return c.Redirect(http.StatusFound, returnTo)
		}
		// If invalid, just return JSON success or redirect to default
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
		"token":   token,
	})
}
