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

	t.Run("successfully renders index", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := NewHTMLHandler()

		// For the test, we can use the actual NewTemplateRenderer since it uses embed
		e.Renderer = NewTemplateRenderer()

		// Note: The embed directive in html_handler.go will embed files relative to that file.
		// Since we are running the test in the same package, it should work fine.

		if assert.NoError(t, h.Index(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			// We can check for a substring since the template is large
			assert.Contains(t, rec.Body.String(), "OAuth2 Provider")
			assert.Contains(t, rec.Body.String(), "Available Endpoints")
		}
	})

    t.Run("renderer uses embedded templates", func(t *testing.T) {
        // Verify NewTemplateRenderer loads the templates correctly
        renderer := NewTemplateRenderer()
        assert.NotNil(t, renderer.templates)
        assert.NotNil(t, renderer.templates.Lookup("index.html"))
    })
}
