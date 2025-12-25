package handlers

import (
	"embed"
	"github.com/labstack/echo/v4"
	"html/template"
)

//go:embed templates/home.html
var templateFS embed.FS

type HTMLHandler struct {
	templates *template.Template
}

func NewHTMLHandler() *HTMLHandler {
	// Parse the embedded template
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		// Verify panic usage: acceptable here as app cannot start without templates
		panic(err)
	}
	return &HTMLHandler{templates: tmpl}
}

func (h *HTMLHandler) ShowHome(c echo.Context) error {
	return h.templates.ExecuteTemplate(c.Response().Writer, "home.html", nil)
}
