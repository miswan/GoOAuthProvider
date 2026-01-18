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
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "System Operational") {
		t.Error("expected body to contain 'System Operational'")
	}
	if !strings.Contains(body, "aria-label=\"Status: Operational\"") {
		t.Error("expected body to contain aria-label")
	}
	if !strings.Contains(body, "#4ade80") {
		t.Error("expected body to contain success color")
	}
}
