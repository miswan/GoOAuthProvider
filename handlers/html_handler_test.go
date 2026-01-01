package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHTMLHandler_Index(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHTMLHandler()

	// Assertions
	if assert.NoError(t, h.Index(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "<title>OAuth2 Provider</title>")
		assert.Contains(t, rec.Body.String(), "Operational")
		assert.Contains(t, rec.Body.String(), "Available Endpoints")

		// Check for accessibility attributes
		assert.Contains(t, rec.Body.String(), `lang="en"`)
		assert.Contains(t, rec.Body.String(), `role="status"`)
		assert.Contains(t, rec.Body.String(), `aria-labelledby="endpoints-title"`)

		// Check for specific endpoint listing
		assert.Contains(t, rec.Body.String(), "/authorize")
		assert.Contains(t, rec.Body.String(), "/token")
	}
}
