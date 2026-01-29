package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHTMLHandler_Index(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHTMLHandler()

	// Assertions
	if err := h.Index(c); err != nil {
		t.Errorf("handler error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "System Operational") {
		t.Errorf("response body does not contain 'System Operational'")
	}
}
