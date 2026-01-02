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
		t.Errorf("Welcome() error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Welcome() status code = %v, want %v", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	expectedSubstrings := []string{
		"OAuth2 Provider",
		"Operational",
		"<!DOCTYPE html>",
		"role=\"status\"", // accessibility check
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(body, s) {
			t.Errorf("Welcome() body missing %q", s)
		}
	}
}
