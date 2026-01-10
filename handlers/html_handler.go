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
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OAuth2 Provider Status</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: #f8fafc;
            color: #1e293b;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            margin: 0;
            line-height: 1.5;
        }
        .container {
            background-color: white;
            padding: 2.5rem;
            border-radius: 0.75rem;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
            max-width: 28rem;
            width: 100%;
            text-align: center;
        }
        h1 {
            color: #0f172a;
            font-size: 1.5rem;
            font-weight: 700;
            margin-top: 0;
            margin-bottom: 1rem;
        }
        .status-icon {
            display: block;
            margin: 0 auto 1rem;
        }
        p {
            color: #475569;
            margin-bottom: 1.5rem;
        }
        .meta {
            font-size: 0.875rem;
            color: #94a3b8;
            border-top: 1px solid #e2e8f0;
            padding-top: 1rem;
            margin-top: 1.5rem;
        }
    </style>
</head>
<body>
    <main class="container">
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="status-icon" role="img" aria-label="System active">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
            <polyline points="22 4 12 14.01 9 11.01"></polyline>
        </svg>
        <h1>System Operational</h1>
        <p>The OAuth2 Provider is running and ready to accept requests.</p>
        <div class="meta">
            Service Status: <strong>Active</strong>
        </div>
    </main>
</body>
</html>
`
	return c.HTML(http.StatusOK, html)
}
