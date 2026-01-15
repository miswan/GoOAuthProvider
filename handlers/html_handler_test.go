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
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "System Online") {
		t.Error("expected body to contain 'System Online'")
	}
	if !strings.Contains(body, `aria-label="System Online"`) {
		t.Error("expected aria-label for accessibility")
	}
}

func TestHTMLHandler_Login(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHTMLHandler()

	if err := h.Login(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "<form id=\"loginForm\">") {
		t.Error("expected login form")
	}
	if !strings.Contains(body, `role="alert"`) {
		t.Error("expected role='alert' for accessibility")
	}
}
