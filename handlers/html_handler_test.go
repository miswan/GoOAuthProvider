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
	handler := NewHTMLHandler()

	// Assertions
	if err := handler.Index(c); err != nil {
		t.Fatalf("handler.Index failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected content type text/html, got %s", contentType)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "OAuth2 Provider Status") {
		t.Errorf("expected body to contain 'OAuth2 Provider Status'")
	}

	if !strings.Contains(body, "System Online") {
		t.Errorf("expected body to contain 'System Online' label")
	}
}
