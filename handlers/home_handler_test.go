package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"oauth2-provider/utils"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHomeHandler_Index(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Load templates
	// When running test in 'handlers' pkg, path is ../templates
	tmpl, err := template.ParseGlob("../templates/*.html")
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	e.Renderer = utils.NewTemplateRenderer(tmpl)

	h := NewHomeHandler()

	// Execute
	if err := h.Index(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	// Assertions
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	expectedTitle := "OAuth2 Provider"
	if !strings.Contains(body, expectedTitle) {
		t.Errorf("expected body to contain '%s', got: %s", expectedTitle, body)
	}

	// Accessibility Check (Basic)
	expectedLang := "lang=\"en\""
	if !strings.Contains(body, expectedLang) {
		t.Errorf("expected html tag to include lang attribute")
	}
}
