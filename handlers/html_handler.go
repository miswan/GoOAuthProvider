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
    <title>OAuth2 Provider - System Status</title>
    <meta name="description" content="OAuth2 Provider Service Status Page">
    <style>
        :root {
            --bg-color: #f8fafc;
            --text-color: #334155;
            --heading-color: #1e293b;
            --card-bg: #ffffff;
            --status-bg: #dcfce7;
            --status-text: #166534;
            --shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
        }
        @media (prefers-color-scheme: dark) {
            :root {
                --bg-color: #0f172a;
                --text-color: #cbd5e1;
                --heading-color: #f1f5f9;
                --card-bg: #1e293b;
                --status-bg: #052e16;
                --status-text: #4ade80;
            }
        }
        body {
            font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background-color: var(--bg-color);
            color: var(--text-color);
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            margin: 0;
            line-height: 1.5;
        }
        main {
            background: var(--card-bg);
            padding: 2.5rem;
            border-radius: 1rem;
            box-shadow: var(--shadow);
            max-width: 24rem;
            width: 100%;
            text-align: center;
        }
        h1 {
            color: var(--heading-color);
            margin-top: 0;
            font-size: 1.5rem;
            font-weight: 700;
        }
        .status-badge {
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.5rem 1rem;
            background-color: var(--status-bg);
            color: var(--status-text);
            border-radius: 9999px;
            font-weight: 600;
            font-size: 0.875rem;
            margin-top: 1rem;
        }
        .status-dot {
            width: 0.5rem;
            height: 0.5rem;
            background-color: currentColor;
            border-radius: 50%;
        }
    </style>
</head>
<body>
    <main>
        <h1>OAuth2 Provider</h1>
        <p>The authentication service is running normally.</p>
        <div class="status-badge" role="status">
            <span class="status-dot" aria-hidden="true"></span>
            <span>System Operational</span>
        </div>
    </main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
