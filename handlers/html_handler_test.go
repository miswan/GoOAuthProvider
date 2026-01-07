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
	handler := NewHTMLHandler()

	// Assertions
	if err := handler.Welcome(c); err != nil {
		t.Errorf("handler.Welcome() error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("handler.Welcome() status = %v, want %v", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "OAuth2 Provider") {
		t.Errorf("handler.Welcome() body does not contain 'OAuth2 Provider'")
	}
	if !strings.Contains(body, "System Operational") {
		t.Errorf("handler.Welcome() body does not contain 'System Operational'")
	}
}
