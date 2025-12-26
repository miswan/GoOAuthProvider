package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) ShowLogin(c echo.Context) error {
	next := c.QueryParam("next")
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"Next": next,
	})
}

func (h *HTMLHandler) ShowRegister(c echo.Context) error {
	next := c.QueryParam("next")
	return c.Render(http.StatusOK, "register.html", map[string]interface{}{
		"Next": next,
	})
}
