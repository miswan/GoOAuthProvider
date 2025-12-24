package handlers

import (
	"embed"
	"github.com/labstack/echo/v4"
	"html/template"
)

//go:embed templates/*.html
var templates embed.FS

type HTMLHandler struct {
	t *template.Template
}

func NewHTMLHandler() *HTMLHandler {
	t := template.Must(template.ParseFS(templates, "templates/*.html"))
	return &HTMLHandler{t: t}
}

func (h *HTMLHandler) ShowHome(c echo.Context) error {
	return h.t.ExecuteTemplate(c.Response().Writer, "home.html", nil)
}

func (h *HTMLHandler) ShowLogin(c echo.Context) error {
	return h.t.ExecuteTemplate(c.Response().Writer, "login.html", nil)
}
