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
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OAuth2 Provider Status</title>
    <style>
        :root {
            --primary: #2563eb;
            --text: #1f2937;
            --bg: #f3f4f6;
            --card-bg: #ffffff;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: var(--bg);
            color: var(--text);
            line-height: 1.5;
            margin: 0;
            padding: 2rem;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
        }
        main {
            background: var(--card-bg);
            padding: 2rem;
            border-radius: 8px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
            max-width: 600px;
            width: 100%;
        }
        h1 {
            color: var(--primary);
            margin-top: 0;
        }
        .status-indicator {
            display: inline-block;
            width: 10px;
            height: 10px;
            background-color: #10b981;
            border-radius: 50%;
            margin-right: 8px;
        }
        ul {
            list-style: none;
            padding: 0;
        }
        li {
            padding: 0.5rem 0;
            border-bottom: 1px solid #e5e7eb;
        }
        li:last-child {
            border-bottom: none;
        }
        code {
            background: #f1f5f9;
            padding: 0.2rem 0.4rem;
            border-radius: 4px;
            font-size: 0.875rem;
            color: #d946ef;
        }
    </style>
</head>
<body>
    <main>
        <h1><span class="status-indicator" aria-label="System Online"></span>OAuth2 Provider</h1>
        <p>The service is running correctly.</p>

        <h2>Available Endpoints</h2>
        <ul>
            <li><code>GET /authorize</code> - Authorization endpoint</li>
            <li><code>POST /token</code> - Token endpoint</li>
            <li><code>GET /userinfo</code> - User info endpoint</li>
            <li><code>POST /register</code> - User registration</li>
            <li><code>POST /login</code> - User login</li>
        </ul>
    </main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
