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
		t.Fatalf("handler error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "System Operational") {
		t.Error("expected body to contain 'System Operational'")
	}
	if !strings.Contains(body, "Log In") {
		t.Error("expected body to contain 'Log In' link")
	}
}

func TestHTMLHandler_Login(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NewHTMLHandler()

	if err := h.Login(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "<form id=\"loginForm\"") {
		t.Error("expected body to contain login form")
	}
	if !strings.Contains(body, "name=\"username\"") {
		t.Error("expected body to contain username field")
	}
	if !strings.Contains(body, "name=\"password\"") {
		t.Error("expected body to contain password field")
	}
}
