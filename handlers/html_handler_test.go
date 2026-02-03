package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestLoginOpenRedirect(t *testing.T) {
	handler := NewHTMLHandler()
	e := echo.New()

	// Malicious URL (Absolute)
	req := httptest.NewRequest(http.MethodGet, "/login?return_to=http://evil.com", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.Login(c); err != nil {
		t.Errorf("Login returned error: %v", err)
	}

	body := rec.Body.String()
	if strings.Contains(body, `value="http://evil.com"`) {
		t.Errorf("Open redirect vulnerability: malicious URL found in hidden input")
	}
	if !strings.Contains(body, `value="/"`) {
		t.Errorf("Expected return_to to be sanitized to /")
	}

	// Malicious URL (Protocol Relative)
	req = httptest.NewRequest(http.MethodGet, "/login?return_to=//evil.com", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	handler.Login(c)
	body = rec.Body.String()
	if strings.Contains(body, `value="//evil.com"`) {
		t.Errorf("Open redirect vulnerability: malicious URL found in hidden input")
	}

	// Valid URL
	req = httptest.NewRequest(http.MethodGet, "/login?return_to=/dashboard", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	handler.Login(c)
	body = rec.Body.String()
	if !strings.Contains(body, `value="/dashboard"`) {
		t.Errorf("Valid relative URL should be preserved")
	}
}
