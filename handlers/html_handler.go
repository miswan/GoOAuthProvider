package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Index(c echo.Context) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>System Status</title>
    <style>
        :root { --bg: #111; --text: #fff; --success: #4ade80; }
        @media (prefers-color-scheme: light) { :root { --bg: #f9fafb; --text: #111827; } }
        body {
            font-family: system-ui, -apple-system, sans-serif;
            background: var(--bg);
            color: var(--text);
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            margin: 0;
            line-height: 1.5;
        }
        .status { text-align: center; }
        .icon { color: var(--success); width: 48px; height: 48px; margin-bottom: 1rem; }
        h1 { margin: 0; font-size: 1.5rem; font-weight: 600; }
        p { margin: 0.5rem 0 0; opacity: 0.8; }
    </style>
</head>
<body>
    <main class="status">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h1>System Operational</h1>
        <p>OAuth2 Provider is running normally.</p>
    </main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
