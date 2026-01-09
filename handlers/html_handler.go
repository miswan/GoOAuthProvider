package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Welcome(c echo.Context) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OAuth2 Provider Status</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #f4f4f9; color: #333; }
        .container { text-align: center; padding: 2.5rem; background: white; border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.08); max-width: 400px; width: 90%; }
        h1 { margin-top: 0; font-size: 1.5rem; margin-bottom: 1rem; }
        .status { color: #15803d; font-weight: 600; display: inline-flex; align-items: center; justify-content: center; gap: 0.5rem; padding: 0.5rem 1rem; background-color: #f0fdf4; border-radius: 20px; border: 1px solid #dcfce7; }
        .dot { width: 8px; height: 8px; background-color: #15803d; border-radius: 50%; display: block; }
        p { color: #666; line-height: 1.5; margin-top: 1.5rem; margin-bottom: 0; }
    </style>
</head>
<body>
    <main class="container">
        <h1>OAuth2 Provider</h1>
        <div class="status" role="status" aria-label="System Status: Online">
            <span class="dot" aria-hidden="true"></span>
            System Online
        </div>
        <p>The OAuth2 service is operational and ready to process requests.</p>
    </main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
