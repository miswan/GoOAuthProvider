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
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()
	expectedContent := []string{
		"System Operational",
		"OAuth2 Provider Status",
		"aria-label=\"System active\"",
	}

	for _, content := range expectedContent {
		if !strings.Contains(body, content) {
			t.Errorf("expected response to contain %q", content)
		}
	}
}
