package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
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
        :root {
            --bg-color: #f9fafb;
            --card-bg: #ffffff;
            --text-primary: #111827;
            --text-secondary: #6b7280;
            --accent-color: #4f46e5;
            --success-color: #10b981;
            --border-color: #e5e7eb;
        }
        @media (prefers-color-scheme: dark) {
            :root {
                --bg-color: #111827;
                --card-bg: #1f2937;
                --text-primary: #f9fafb;
                --text-secondary: #9ca3af;
                --border-color: #374151;
            }
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: var(--bg-color);
            color: var(--text-primary);
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            margin: 0;
            line-height: 1.5;
        }
        .card {
            background-color: var(--card-bg);
            padding: 2.5rem;
            border-radius: 1rem;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
            max-width: 24rem;
            width: 100%;
            border: 1px solid var(--border-color);
            text-align: center;
        }
        .status-badge {
            display: inline-flex;
            align-items: center;
            padding: 0.25rem 0.75rem;
            background-color: rgba(16, 185, 129, 0.1);
            color: var(--success-color);
            border-radius: 9999px;
            font-size: 0.875rem;
            font-weight: 500;
            margin-bottom: 1.5rem;
        }
        .status-dot {
            width: 0.5rem;
            height: 0.5rem;
            background-color: currentColor;
            border-radius: 50%;
            margin-right: 0.5rem;
        }
        h1 {
            font-size: 1.5rem;
            font-weight: 700;
            margin: 0 0 0.5rem 0;
            color: var(--text-primary);
        }
        p {
            color: var(--text-secondary);
            margin: 0 0 2rem 0;
        }
        .footer {
            font-size: 0.875rem;
            color: var(--text-secondary);
            border-top: 1px solid var(--border-color);
            padding-top: 1.5rem;
            margin-top: 1.5rem;
        }
    </style>
</head>
<body>
    <main class="card">
        <div class="status-badge" role="status" aria-label="System status: Operational">
            <span class="status-dot" aria-hidden="true"></span>
            Operational
        </div>
        <h1>OAuth2 Provider</h1>
        <p>The authentication service is running and ready to handle requests.</p>
        <div class="footer">
            Managed by Palette 🎨
        </div>
    </main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
