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
    <title>OAuth2 Provider</title>
    <style>
        :root {
            --primary: #2563eb;
            --success: #16a34a;
            --text: #1f2937;
            --bg: #f3f4f6;
            --card-bg: #ffffff;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: var(--bg);
            color: var(--text);
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            margin: 0;
            padding: 20px;
        }
        .card {
            background: var(--card-bg);
            padding: 2rem;
            border-radius: 12px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
            max-width: 400px;
            width: 100%;
        }
        h1 {
            margin-top: 0;
            font-size: 1.5rem;
            margin-bottom: 0.5rem;
        }
        .status {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            font-weight: 500;
            margin-bottom: 1.5rem;
            color: var(--success);
            font-size: 0.875rem;
        }
        .status-dot {
            width: 8px;
            height: 8px;
            background-color: var(--success);
            border-radius: 50%;
            display: inline-block;
        }
        .info {
            line-height: 1.5;
            color: #4b5563;
            margin-bottom: 1.5rem;
        }
        .endpoints {
            background: #f9fafb;
            padding: 1rem;
            border-radius: 8px;
            font-size: 0.875rem;
        }
        .endpoint-item {
            display: flex;
            justify-content: space-between;
            margin-bottom: 0.5rem;
        }
        .endpoint-item:last-child {
            margin-bottom: 0;
        }
        .method {
            font-weight: bold;
            color: var(--primary);
        }
        .path {
            font-family: monospace;
            color: #374151;
        }
    </style>
</head>
<body>
    <div class="card">
        <h1>OAuth2 Provider</h1>
        <div class="status" role="status">
            <span class="status-dot" aria-hidden="true"></span>
            Operational
        </div>
        <p class="info">
            Welcome to the OAuth2 Provider service. This API handles user authentication and authorization using the OAuth 2.0 protocol.
        </p>
        <div class="endpoints">
            <div class="endpoint-item">
                <span class="method">GET</span>
                <span class="path">/authorize</span>
            </div>
            <div class="endpoint-item">
                <span class="method">POST</span>
                <span class="path">/token</span>
            </div>
             <div class="endpoint-item">
                <span class="method">POST</span>
                <span class="path">/register</span>
            </div>
        </div>
    </div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
