package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHTMLHandler_Welcome(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHTMLHandler()

	// Execute
	if err := h.Welcome(c); err != nil {
		t.Fatalf("Welcome failed: %v", err)
	}

	// Verify
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()
	expectedSubstrings := []string{
		"OAuth2 Provider Status",
		"Systems Operational",
		"role=\"status\"",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(body, s) {
			t.Errorf("Expected body to contain %q", s)
		}
	}
}
