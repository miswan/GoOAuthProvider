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
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OAuth2 Provider - Operational</title>
    <style>
        :root { --bg: #fff; --text: #1f2937; --success: #10b981; --subtle: #6b7280; }
        @media (prefers-color-scheme: dark) { :root { --bg: #111; --text: #f3f4f6; --success: #34d399; --subtle: #9ca3af; } }
        body { font-family: system-ui, -apple-system, sans-serif; background: var(--bg); color: var(--text); display: grid; place-items: center; min-height: 100vh; margin: 0; line-height: 1.5; }
        .status-card { text-align: center; padding: 2rem; animation: fade-in 0.5s ease-out; }
        .icon { width: 64px; height: 64px; color: var(--success); margin-bottom: 1rem; }
        h1 { margin: 0 0 0.5rem; font-size: 1.5rem; }
        p { margin: 0; color: var(--subtle); }
        @keyframes fade-in { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
    </style>
</head>
<body>
    <main class="status-card">
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
            <polyline points="22 4 12 14.01 9 11.01"></polyline>
        </svg>
        <h1>System Operational</h1>
        <p>OAuth2 Provider is running successfully.</p>
    </main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
