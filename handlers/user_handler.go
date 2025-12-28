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

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	req := new(models.UserLogin)
	// Bind JSON or Form data
	if err := c.Bind(req); err != nil {
		// Fallback for form data if strict binding fails or isn't used
		req.Username = c.FormValue("username")
		req.Password = c.FormValue("password")
	}

	if req.Username == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username and password required")
	}

	user, err := h.userService.Login(req)
	if err != nil {
		redirectTo := c.FormValue("redirect_to")
		if redirectTo != "" {
			return c.Render(http.StatusUnauthorized, "login.html", map[string]interface{}{
				"Error":      "Invalid credentials",
				"RedirectTo": redirectTo,
			})
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// Generate JWT for the session
	token, err := utils.GenerateJWT(user.ID, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	// Set session cookie
	c.SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600 * 24, // 24 hours
	})

	// Check if we need to redirect
	redirectTo := c.FormValue("redirect_to")
	if redirectTo != "" {
		return c.Redirect(http.StatusFound, redirectTo)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
	})
}
