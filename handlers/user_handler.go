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

func (h *UserHandler) LoginPage(c echo.Context) error {
	returnTo := c.QueryParam("return_to")
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"ReturnTo": returnTo,
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	// Check if this is a form submission or JSON request
	contentType := c.Request().Header.Get("Content-Type")
	isJSON := contentType == "application/json"

	req := new(models.UserLogin)
	if err := c.Bind(req); err != nil {
		if isJSON {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error": "Invalid request",
		})
	}

	user, err := h.userService.Login(req)
	if err != nil {
		if isJSON {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error": "Invalid credentials",
		})
	}

	// Generate session token (using JWT for simplicity, but stored in cookie)
	token, err := services.GenerateSessionToken(user.ID)
	if err != nil {
		if isJSON {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error": "Internal server error",
		})
	}

	// Set cookie
	cookie := new(http.Cookie)
	cookie.Name = "session_token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	if isJSON {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Login successful",
			"user_id": strconv.FormatUint(uint64(user.ID), 10),
			"token":   token,
		})
	}

	// Redirect to return_to if present
	returnTo := c.FormValue("return_to")
	if returnTo != "" {
		// Validate that return_to starts with / and not // (protocol-relative)
		if len(returnTo) >= 1 && returnTo[0] == '/' {
			if len(returnTo) == 1 || returnTo[1] != '/' {
				return c.Redirect(http.StatusFound, returnTo)
			}
		}
	}

	return c.Redirect(http.StatusFound, "/")
}