package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHTMLHandler_Home(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHTMLHandler()

	// Assertions
	if err := h.Home(c); err != nil {
		t.Errorf("handler error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected content type text/html, got %s", contentType)
	}

	body := rec.Body.String()
	expectedStrings := []string{
		"<!DOCTYPE html>",
		"<title>OAuth2 Provider - System Status</title>",
		"System Operational",
		"role=\"status\"",
		"POST /register",
	}

	for _, s := range expectedStrings {
		if !strings.Contains(body, s) {
			t.Errorf("response body missing expected string: %s", s)
		}
	}
}
