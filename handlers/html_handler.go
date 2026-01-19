package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Index(c echo.Context) error {
	const html = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OAuth2 Provider</title>
    <style>
        :root { --bg: #111; --text: #eee; --success: #4ade80; }
        body { background: var(--bg); color: var(--text); font-family: system-ui, -apple-system, sans-serif; display: grid; place-items: center; height: 100vh; margin: 0; }
        .status { display: flex; align-items: center; gap: 0.75rem; font-weight: 500; font-size: 1.1rem; }
        .dot { width: 10px; height: 10px; background: var(--success); border-radius: 50%; box-shadow: 0 0 12px var(--success); animation: pulse 2s infinite; }
        @keyframes pulse { 0% { opacity: 1; } 50% { opacity: 0.5; } 100% { opacity: 1; } }
        @media (prefers-reduced-motion: reduce) { .dot { animation: none; } }
    </style>
</head>
<body>
    <div class="status" role="status" aria-label="System status: Operational">
        <div class="dot" aria-hidden="true"></div>
        <span>System Operational</span>
    </div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
