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
	// Try binding JSON first, then form data
	req := new(models.UserLogin)
	if err := c.Bind(req); err != nil {
		// If binding fails, try to parse form values manually as fallback
		req.Username = c.FormValue("username")
		req.Password = c.FormValue("password")
		if req.Username == "" || req.Password == "" {
			// If both JSON and Form binding fail (or produce empty values)
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid login credentials")
		}
	} else if req.Username == "" || req.Password == "" {
        // If Bind passed but values are empty (e.g. empty JSON or form fields handled by Bind but missing validation)
        // Try fallback to explicit form values if not already set
        if u := c.FormValue("username"); u != "" {
            req.Username = u
        }
        if p := c.FormValue("password"); p != "" {
            req.Password = p
        }

        if req.Username == "" || req.Password == "" {
            return echo.NewHTTPError(http.StatusBadRequest, "Invalid login credentials")
        }
    }

	user, err := h.userService.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// Generate session token (using JWT for simplicity in this cookie)
	// In a real app, this might be a random session ID stored in DB/Redis
	// But we'll reuse the utils.GenerateJWT for now
	token, err := services.GenerateSessionToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate session")
	}

	// Set cookie
	c.SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
	})

	redirectTo := c.FormValue("redirect_to")
	if redirectTo != "" {
		return c.Redirect(http.StatusFound, redirectTo)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
	})
}