package handlers

import (
	"net/http"
	"net/http/httptest"
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
		t.Fatalf("Welcome handler returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()
	expectedStrings := []string{
		"OAuth2 Provider Status",
		"System Online",
		"role=\"status\"",
	}

	for _, expected := range expectedStrings {
		if !contains(body, expected) {
			t.Errorf("Response body missing expected content: %s", expected)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && search(s, substr)
}

func search(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
