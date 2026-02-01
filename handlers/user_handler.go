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

	// Attempt to bind from JSON or Form
	if err := c.Bind(req); err != nil {
		// If bind fails or is partial, we ensure we check form values manually if needed
		// But c.Bind usually handles Form values if Content-Type is application/x-www-form-urlencoded
	}

	// Fallback/Ensure values if Bind didn't work as expected for mixed content types
	if req.Username == "" {
		req.Username = c.FormValue("username")
	}
	if req.Password == "" {
		req.Password = c.FormValue("password")
	}

	if req.Username == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username and password are required")
	}

	user, err := h.userService.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// Generate Token
	token, err := utils.GenerateJWT(user.ID, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	// Set Cookie
	c.SetCookie(&http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Should be true in production
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// Check for return_to parameter (from form or query)
	returnTo := c.FormValue("return_to")
	if returnTo == "" {
		returnTo = c.QueryParam("return_to")
	}

	if returnTo != "" {
		return c.Redirect(http.StatusFound, returnTo)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user_id": strconv.FormatUint(uint64(user.ID), 10),
		"token":   token,
	})
}
