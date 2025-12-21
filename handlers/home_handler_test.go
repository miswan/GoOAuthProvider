package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"

	"github.com/labstack/echo/v4"
)

// MockTemplateRenderer implements echo.Renderer interface for testing
type MockTemplateRenderer struct{}

func (t *MockTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl := template.Must(template.New(name).Parse("<html><body><h1>OAuth2 Provider</h1></body></html>"))
	return tmpl.Execute(w, data)
}

func TestHomeHandler_Index(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Use a mock renderer because we don't want to depend on actual template files in unit tests
	// or we can point to the actual templates. For unit tests, mocking is safer.
	// But since we want to test that the template exists, let's try to use the real one if possible,
	// but the path might be tricky. Let's just test the handler logic.

	h := NewHomeHandler()

	// We need to set a renderer on Echo to test c.Render
	// However, c.Render calls e.Renderer.Render.
	// We can mock the renderer to verify it's called with "index.html".

	mockRenderer := &testRenderer{}
	e.Renderer = mockRenderer

	// Execute
	if err := h.Index(c); err != nil {
		t.Fatalf("handler execution failed: %v", err)
	}

	// Verify
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	if mockRenderer.templateName != "index.html" {
		t.Errorf("expected template name 'index.html', got '%s'", mockRenderer.templateName)
	}
}

type testRenderer struct {
	templateName string
}

func (t *testRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	t.templateName = name
	return nil
}
