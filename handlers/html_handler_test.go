package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/labstack/echo/v4"
)

func TestIndexHandler(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHTMLHandler()

	// Assertions
	if err := h.Index(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rec.Code)
	}

	// Basic content check
	body := rec.Body.String()
	if body == "" {
		t.Error("Expected body content, got empty string")
	}

	expectedTitle := "<title>System Status</title>"
	if !strings.Contains(body, expectedTitle) {
		t.Errorf("Expected body to contain %q", expectedTitle)
	}
}
