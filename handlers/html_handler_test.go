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
		t.Fatalf("Index() error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()

	// Accessibility checks
	if !strings.Contains(body, `lang="en"`) {
		t.Error("missing lang attribute")
	}
	if !strings.Contains(body, `role="main"`) {
		t.Error("missing main role")
	}
	if !strings.Contains(body, `role="status"`) {
		t.Error("missing status role")
	}

	// Content checks
	if !strings.Contains(body, "OAuth2 Provider") {
		t.Error("missing title")
	}
}
