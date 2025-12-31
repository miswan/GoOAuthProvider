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
		t.Errorf("handler error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "OAuth2 Provider") {
		t.Errorf("expected body to contain 'OAuth2 Provider', got %s", body)
	}
	if !strings.Contains(body, "Operational") {
		t.Errorf("expected body to contain 'Operational', got %s", body)
	}
	if !strings.Contains(body, "<html lang=\"en\">") {
		t.Errorf("expected body to contain lang attribute, got %s", body)
	}
}
