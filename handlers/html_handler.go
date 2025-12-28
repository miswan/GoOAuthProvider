package handlers

import (
	"html/template"
	"io"
	"net/http"
	"embed"
	"github.com/labstack/echo/v4"
)

//go:embed templates/*.html
var templatesFS embed.FS

type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{
		templates: template.Must(template.ParseFS(templatesFS, "templates/*.html")),
	}
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Login(c echo.Context) error {
	redirectTo := c.QueryParam("redirect_to")
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"RedirectTo": redirectTo,
	})
}
