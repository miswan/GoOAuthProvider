package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestServeHome(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHTMLHandler()

	// Assertions
	if err := h.ServeHome(c); err != nil {
		t.Errorf("ServeHome() error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected Content-Type to contain text/html, got %s", contentType)
	}

	body := rec.Body.String()
	expectedSnippet := "OAuth2 Provider Status"
	if !strings.Contains(body, expectedSnippet) {
		t.Errorf("expected body to contain %q, got %q", expectedSnippet, body)
	}
}
