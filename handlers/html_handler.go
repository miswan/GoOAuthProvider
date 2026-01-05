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
    <title>OAuth2 Provider Status</title>
    <style>
        :root {
            --primary: #2563eb;
            --text: #1f2937;
            --bg: #f3f4f6;
            --card: #ffffff;
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
            line-height: 1.5;
        }
        main {
            background: var(--card);
            padding: 2rem;
            border-radius: 8px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
            max-width: 400px;
            width: 100%;
            text-align: center;
        }
        h1 { margin-top: 0; color: var(--primary); }
        .status {
            display: inline-block;
            padding: 0.25rem 0.75rem;
            border-radius: 9999px;
            background-color: #dcfce7;
            color: #166534;
            font-weight: 600;
            font-size: 0.875rem;
            margin-bottom: 1.5rem;
        }
        p { margin-bottom: 1.5rem; color: #4b5563; }
        .footer { margin-top: 2rem; font-size: 0.875rem; color: #6b7280; }
    </style>
</head>
<body>
    <main>
        <h1>OAuth2 Provider</h1>
        <div class="status" role="status">● Systems Operational</div>
        <p>The OAuth2 authentication service is running and ready to accept requests.</p>
        <div class="footer">
            Version 1.0.0
        </div>
    </main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
