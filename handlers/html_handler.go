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
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OAuth2 Provider</title>
    <style>
        :root {
            --primary: #2563eb;
            --bg: #f8fafc;
            --text: #1e293b;
            --card-bg: #ffffff;
            --code-bg: #f1f5f9;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            line-height: 1.6;
            max-width: 800px;
            margin: 0 auto;
            padding: 2rem;
            background: var(--bg);
            color: var(--text);
        }
        main {
            background: var(--card-bg);
            padding: 2.5rem;
            border-radius: 12px;
            box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1);
        }
        h1 {
            color: var(--primary);
            margin-top: 0;
            font-size: 2rem;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }
        .status-badge {
            display: inline-flex;
            align-items: center;
            padding: 0.25rem 0.75rem;
            background: #dcfce7;
            color: #166534;
            border-radius: 9999px;
            font-size: 0.875rem;
            font-weight: 500;
        }
        .endpoint-list {
            list-style: none;
            padding: 0;
            margin-top: 2rem;
        }
        .endpoint-item {
            padding: 1rem;
            border: 1px solid #e2e8f0;
            border-radius: 8px;
            margin-bottom: 1rem;
        }
        .method {
            font-weight: bold;
            color: var(--primary);
            margin-right: 0.5rem;
        }
        code {
            background: var(--code-bg);
            padding: 0.2rem 0.4rem;
            border-radius: 4px;
            font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
            font-size: 0.9em;
        }
        .footer {
            margin-top: 2rem;
            text-align: center;
            font-size: 0.875rem;
            color: #64748b;
        }
    </style>
</head>
<body>
    <main>
        <h1>
            <span>OAuth2 Provider</span>
            <span class="status-badge" role="status">● System Operational</span>
        </h1>
        <p>
            Welcome to the OAuth2 Identity Provider service. This service handles user authentication and authorization.
        </p>

        <h2>Available Endpoints</h2>
        <ul class="endpoint-list">
            <li class="endpoint-item">
                <div><span class="method">GET</span> <code>/authorize</code></div>
                <p>Initiate the OAuth2 authorization flow.</p>
            </li>
            <li class="endpoint-item">
                <div><span class="method">POST</span> <code>/token</code></div>
                <p>Exchange authorization code for access token.</p>
            </li>
            <li class="endpoint-item">
                <div><span class="method">GET</span> <code>/userinfo</code></div>
                <p>Retrieve authenticated user information.</p>
            </li>
        </ul>
    </main>
    <div class="footer">
        <p>Managed by Palette 🎨</p>
    </div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
