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

	// Execution
	if err := h.Index(c); err != nil {
		t.Fatalf("handler execution failed: %v", err)
	}

	// Assertions
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "System Operational") {
		t.Error("body does not contain 'System Operational'")
	}
	if !strings.Contains(body, "<html lang=\"en\">") {
		t.Error("body does not contain lang attribute")
	}
}
