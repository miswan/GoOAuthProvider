package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}
