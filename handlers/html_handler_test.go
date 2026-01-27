package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHTMLHandler_Index(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NewHTMLHandler()

	if err := h.Index(c); err != nil {
		t.Fatalf("handler execution failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expectedContent := "System Operational"
	if !strings.Contains(rec.Body.String(), expectedContent) {
		t.Errorf("expected body to contain %q, got %q", expectedContent, rec.Body.String())
	}

	// Verify accessibility
	if !strings.Contains(rec.Body.String(), `role="status"`) {
		t.Error("missing role=\"status\"")
	}
}
