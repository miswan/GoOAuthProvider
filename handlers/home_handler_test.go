package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"strings"
)

// Mock TemplateRenderer for testing
type MockTemplateRenderer struct {
	templates *template.Template
}

func (t *MockTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func TestHome(t *testing.T) {
	// Setup
	e := echo.New()

	// Load the actual template to ensure it parses and renders correctly
	renderer := &MockTemplateRenderer{
		templates: template.Must(template.ParseGlob("../templates/*.html")),
	}
	e.Renderer = renderer

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if err := Home(c); err != nil {
		t.Errorf("Home() error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Palette OAuth2 Provider") {
		t.Error("body does not contain 'Palette OAuth2 Provider'")
	}
	if !strings.Contains(body, "System Operational") {
		t.Error("body does not contain 'System Operational'")
	}
}
