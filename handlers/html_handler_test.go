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

	// Assertions
	if err := h.Welcome(c); err != nil {
		t.Fatalf("Welcome handler returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("Expected Content-Type to contain 'text/html', got '%s'", contentType)
	}

	body := rec.Body.String()
	expectedSubstrings := []string{
		"OAuth2 Provider",
		"System Operational",
		"Painted with 🎨 by Palette",
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(body, s) {
			t.Errorf("Expected response body to contain '%s', but it didn't", s)
		}
	}
}
