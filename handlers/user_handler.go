package handlers

import (
	"html/template"
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

func (h *UserHandler) LoginView(c echo.Context) error {
	return_to := c.QueryParam("return_to")
	// Simple HTML login form
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Login</title>
	</head>
	<body>
		<h1>Login</h1>
		<form action="/login" method="POST">
			<input type="hidden" name="return_to" value="{{.ReturnTo}}">
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
	</body>
	</html>
	`
	t, err := template.New("login").Parse(tmpl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	data := map[string]string{
		"ReturnTo": return_to,
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return t.Execute(c.Response().Writer, data)
}

func (h *UserHandler) Login(c echo.Context) error {
	req := new(models.UserLogin)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Support form data binding if json binding fails or isn't used
	if req.Username == "" {
		req.Username = c.FormValue("username")
		req.Password = c.FormValue("password")
	}

	user, err := h.userService.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// Generate session token
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
	c.SetCookie(cookie)

	// Check for return_to
	returnTo := c.FormValue("return_to")
	if returnTo != "" {
		// Basic open redirect protection: ensure it's relative or same domain
		// Ensure it starts with / but not // (protocol relative)
		if strings.HasPrefix(returnTo, "/") && !strings.HasPrefix(returnTo, "//") {
			return c.Redirect(http.StatusFound, returnTo)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
	})
}