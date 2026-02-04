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
		t.Fatalf("Index() error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Index() status = %v, want %v", rec.Code, http.StatusOK)
	}

	if !strings.Contains(rec.Body.String(), "System Operational") {
		t.Errorf("Index() body does not contain 'System Operational'")
	}
}
