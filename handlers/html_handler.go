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
    <title>OAuth2 Provider</title>
    <style>
        :root {
            --bg: #111;
            --text: #eee;
            --status: #4ade80;
        }
        @media (prefers-color-scheme: light) {
            :root { --bg: #f9fafb; --text: #111; }
        }
        body {
            margin: 0;
            font-family: system-ui, -apple-system, sans-serif;
            background: var(--bg);
            color: var(--text);
            display: flex;
            align-items: center;
            justify-content: center;
            height: 100vh;
        }
        main { text-align: center; }
        h1 { margin: 0 0 1rem; font-weight: 600; letter-spacing: -0.025em; }
        .badge {
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.5rem 1rem;
            background: rgba(74, 222, 128, 0.1);
            color: var(--status);
            border-radius: 9999px;
            font-size: 0.875rem;
            font-weight: 500;
        }
        .dot {
            width: 0.5rem;
            height: 0.5rem;
            background: currentColor;
            border-radius: 50%;
            position: relative;
        }
        .dot::after {
            content: '';
            position: absolute;
            inset: -2px;
            border-radius: 50%;
            background: currentColor;
            opacity: 0.4;
            animation: pulse 2s infinite;
        }
        @keyframes pulse {
            0% { transform: scale(1); opacity: 0.4; }
            50% { transform: scale(2); opacity: 0; }
            100% { transform: scale(1); opacity: 0; }
        }
        @media (prefers-reduced-motion: reduce) {
            .dot::after { animation: none; }
        }
    </style>
</head>
<body>
    <main>
        <h1>OAuth2 Provider</h1>
        <div class="badge">
            <span class="dot" aria-hidden="true"></span>
            System Operational
        </div>
    </main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
