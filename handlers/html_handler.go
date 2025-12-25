package handlers

import (
	"embed"
	"github.com/labstack/echo/v4"
	"net/http"
	"text/template"
)

//go:embed templates/*.html
var templates embed.FS

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) ShowLogin(c echo.Context) error {
	tmpl, err := template.ParseFS(templates, "templates/login.html")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Template error: "+err.Error())
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return tmpl.Execute(c.Response().Writer, nil)
}
