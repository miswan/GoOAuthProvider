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
    <title>OAuth2 Provider - Status</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; line-height: 1.5; max-width: 40rem; margin: 2rem auto; padding: 0 1rem; color: #1f2937; background-color: #f9fafb; }
        main { background: white; padding: 2rem; border-radius: 0.5rem; box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1); }
        h1 { color: #111827; margin-top: 0; font-size: 1.875rem; }
        .status { margin-top: 1.5rem; padding: 1rem; background: #eff6ff; border-radius: 0.375rem; border-left: 4px solid #3b82f6; }
        .status-dot { color: #10b981; margin-right: 0.5rem; }
        .footer { margin-top: 2rem; font-size: 0.875rem; color: #6b7280; text-align: center; }
    </style>
</head>
<body>
    <main>
        <h1>OAuth2 Provider</h1>
        <p>Welcome to the OAuth2 Provider service.</p>

        <div class="status" role="status" aria-live="polite">
            <p><strong>System Status:</strong> <span class="status-dot">●</span> Operational</p>
        </div>

        <p style="margin-top: 1.5rem;">
            This service provides OAuth2 authentication and authorization endpoints.
        </p>
    </main>
    <footer class="footer">
        <p>Powered by Go and Echo</p>
    </footer>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
