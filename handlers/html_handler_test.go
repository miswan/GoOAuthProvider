package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHTMLHandler_Welcome(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NewHTMLHandler()

	if err := h.Welcome(c); err != nil {
		t.Fatalf("Welcome handler failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "<html lang=\"en\">") {
		t.Error("response missing lang attribute")
	}
	if !strings.Contains(body, "role=\"status\"") {
		t.Error("response missing status role")
	}
	if !strings.Contains(body, "OAuth2 Provider") {
		t.Error("response missing title")
	}
}
