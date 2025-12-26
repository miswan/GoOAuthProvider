package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) ServeHome(c echo.Context) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OAuth2 Provider Status</title>
    <style>
        :root {
            --bg-color: #f9fafb;
            --card-bg: #ffffff;
            --text-primary: #111827;
            --text-secondary: #6b7280;
            --success-color: #10b981;
            --border-color: #e5e7eb;
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
            padding: 2rem;
            border-radius: 0.5rem;
            box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
            max-width: 24rem;
            width: 100%;
            text-align: center;
            border: 1px solid var(--border-color);
        }
        .status-indicator {
            display: inline-flex;
            align-items: center;
            background-color: #d1fae5;
            color: #065f46;
            padding: 0.25rem 0.75rem;
            border-radius: 9999px;
            font-size: 0.875rem;
            font-weight: 500;
            margin-bottom: 1.5rem;
        }
        .dot {
            height: 0.5rem;
            width: 0.5rem;
            background-color: var(--success-color);
            border-radius: 50%;
            margin-right: 0.5rem;
        }
        h1 {
            font-size: 1.5rem;
            font-weight: 700;
            margin: 0 0 0.5rem 0;
        }
        p {
            color: var(--text-secondary);
            margin: 0 0 1.5rem 0;
        }
        .link {
            color: #2563eb;
            text-decoration: none;
            font-weight: 500;
        }
        .link:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <main class="card">
        <div class="status-indicator" role="status">
            <span class="dot" aria-hidden="true"></span>
            System Operational
        </div>
        <h1>OAuth2 Provider</h1>
        <p>The authentication service is running and ready to handle requests.</p>
        <a href="https://oauth.net/2/" class="link" target="_blank" rel="noopener noreferrer">Learn about OAuth 2.0 &rarr;</a>
    </main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
