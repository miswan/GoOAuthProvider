package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHTMLHandler_Welcome(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NewHTMLHandler()

	// Assertions
	if assert.NoError(t, h.Welcome(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "<!DOCTYPE html>")
		assert.Contains(t, rec.Body.String(), "OAuth2 Provider")
		assert.Contains(t, rec.Body.String(), "System Status")
	}
}
