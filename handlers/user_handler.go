package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"oauth2-provider/models"
	"oauth2-provider/services"
	"oauth2-provider/utils"
	"strconv"
	"time"
	"html/template"
	"bytes"
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
	// Try binding JSON first
	if err := c.Bind(req); err != nil {
		// If JSON binding fails, try form values if provided (for browser login)
		req.Username = c.FormValue("username")
		req.Password = c.FormValue("password")
		if req.Username == "" || req.Password == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid login request")
		}
	}

	user, err := h.userService.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// Generate session token (JWT)
	token, err := utils.GenerateJWT(user.ID, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate session")
	}

	// Set cookie
	cookie := new(http.Cookie)
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	// Check if there is a redirect_to parameter
	redirectTo := c.QueryParam("redirect_to")
	if redirectTo != "" {
		// Basic validation to prevent open redirect vulnerabilities
		// In production, validate domain whitelist or relative path
		if utils.IsValidRedirect(c.Request().Host, redirectTo) {
			return c.Redirect(http.StatusFound, redirectTo)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
	})
}

// ShowLogin serves the HTML login page
func (h *UserHandler) ShowLogin(c echo.Context) error {
	redirectTo := c.QueryParam("redirect_to")

	// Use html/template to prevent XSS
	tmplStr := `
<!DOCTYPE html>
<html>
<head>
    <title>Login - OAuth Provider</title>
	<style>
		body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; }
		.login-container { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); width: 300px; }
		.form-group { margin-bottom: 1rem; }
		label { display: block; margin-bottom: 0.5rem; }
		input { width: 100%; padding: 0.5rem; border: 1px solid #ccc; border-radius: 4px; box-sizing: border-box; }
		button { width: 100%; padding: 0.75rem; background-color: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
		button:hover { background-color: #0056b3; }
	</style>
</head>
<body>
    <div class="login-container">
        <h2>Login</h2>
        <form action="/login" method="post">
            <input type="hidden" name="redirect_to" value="{{.}}">
            <div class="form-group">
                <label for="username">Username</label>
                <input type="text" id="username" name="username" required>
            </div>
            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" name="password" required>
            </div>
            <button type="submit">Login</button>
        </form>
    </div>
</body>
</html>
`
	t, err := template.New("login").Parse(tmplStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Template error")
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, redirectTo); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Template execution error")
	}

	return c.HTML(http.StatusOK, buf.String())
}
