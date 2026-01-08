package handlers

import (
	"github.com/labstack/echo/v4"
	"html/template"
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

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	// If GET, show login form
	if c.Request().Method == http.MethodGet {
		redirectTo := c.QueryParam("redirect_to")

		// Use html/template to prevent XSS
		tmpl := `
			<form method="POST" action="/login">
				<input type="text" name="username" placeholder="Username" required>
				<input type="password" name="password" placeholder="Password" required>
				<input type="hidden" name="redirect_to" value="{{.RedirectTo}}">
				<button type="submit">Login</button>
			</form>
		`
		t, err := template.New("login").Parse(tmpl)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Template error")
		}

		data := struct {
			RedirectTo string
		}{
			RedirectTo: redirectTo,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return t.Execute(c.Response().Writer, data)
	}

	req := new(models.UserLogin)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userService.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// Generate Session Token
	token, err := utils.GenerateJWT(user.ID, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate session")
	}

	// Set Cookie
	cookie := new(http.Cookie)
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	if req.RedirectTo != "" {
		// Validate RedirectTo to prevent open redirect
		if utils.IsValidRedirect(req.RedirectTo) {
			return c.Redirect(http.StatusFound, req.RedirectTo)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
	})
}
