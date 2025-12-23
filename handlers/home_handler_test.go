package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, Home(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "<title>OAuth2 Provider</title>")
		assert.Contains(t, rec.Body.String(), "System is operational")
		// Verify accessibility elements
		assert.Contains(t, rec.Body.String(), `lang="en"`)
		assert.Contains(t, rec.Body.String(), `role="status"`)
	}
}
