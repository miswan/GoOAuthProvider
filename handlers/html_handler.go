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
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OAuth2 Provider</title>
    <style>
        :root {
            --bg: #111;
            --text: #eee;
            --accent: #3b82f6;
            --success: #4ade80;
        }
        @media (prefers-color-scheme: light) {
            :root {
                --bg: #f9fafb;
                --text: #1f2937;
                --accent: #2563eb;
                --success: #16a34a;
            }
        }
        body {
            font-family: system-ui, -apple-system, sans-serif;
            background: var(--bg);
            color: var(--text);
            display: flex;
            align-items: center;
            justify-content: center;
            height: 100vh;
            margin: 0;
        }
        .status {
            text-align: center;
            padding: 2rem;
            border: 1px solid #333;
            border-radius: 8px;
            background: rgba(255, 255, 255, 0.05);
        }
        @media (prefers-color-scheme: light) {
            .status {
                border-color: #e5e7eb;
                background: white;
                box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
            }
        }
        .dot {
            display: inline-block;
            width: 12px;
            height: 12px;
            background: var(--success);
            border-radius: 50%;
            margin-right: 8px;
        }
        h1 { margin: 0 0 0.5rem 0; font-size: 1.5rem; }
        p { margin: 0; opacity: 0.8; }
    </style>
</head>
<body>
    <div class="status">
        <h1><span class="dot"></span>System Operational</h1>
        <p>OAuth2 Provider is running.</p>
    </div>
</body>
</html>
`
	return c.HTML(http.StatusOK, html)
}
