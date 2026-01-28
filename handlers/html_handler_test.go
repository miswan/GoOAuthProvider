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
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "System Operational") {
		t.Errorf("expected body to contain 'System Operational', got %s", body)
	}

	// Accessibility checks (basic)
	if !strings.Contains(body, `lang="en"`) {
		t.Errorf("expected lang attribute")
	}
	if !strings.Contains(body, `aria-hidden="true"`) {
		t.Errorf("expected aria-hidden on icon")
	}
}
