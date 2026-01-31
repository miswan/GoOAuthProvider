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

	// Execute
	if err := h.Index(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	// Verify Status
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Verify Content
	body := rec.Body.String()
	if !strings.Contains(body, "System Operational") {
		t.Errorf("expected body to contain 'System Operational', got %s", body)
	}
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Errorf("expected valid HTML doctype")
	}
}
