package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHTMLHandler_Index(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NewHTMLHandler()

	if err := h.Index(c); err != nil {
		t.Errorf("Index() error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "System Operational") {
		t.Errorf("expected body to contain 'System Operational', got %s", body)
	}
	if !strings.Contains(body, "OAuth2 Provider is running successfully") {
		t.Errorf("expected body to contain 'OAuth2 Provider is running successfully', got %s", body)
	}
}
