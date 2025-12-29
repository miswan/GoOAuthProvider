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
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 650px; margin: 40px auto; padding: 0 20px; }
        h1 { color: #2563eb; }
        .status { padding: 8px 12px; background: #dcfce7; color: #166534; border-radius: 6px; display: inline-flex; align-items: center; font-weight: 500; font-size: 0.9em; }
        .status::before { content: ""; width: 8px; height: 8px; background-color: #166534; border-radius: 50%; margin-right: 8px; }
        .card { border: 1px solid #e5e7eb; border-radius: 8px; padding: 24px; margin-top: 24px; box-shadow: 0 1px 3px rgba(0,0,0,0.05); }
        h2 { font-size: 1.25em; margin-top: 0; margin-bottom: 16px; color: #1f2937; }
        ul { padding-left: 20px; margin: 0; }
        li { margin-bottom: 8px; color: #4b5563; }
        code { background: #f3f4f6; padding: 2px 6px; border-radius: 4px; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; font-size: 0.9em; color: #db2777; }
        .footer { margin-top: 48px; color: #9ca3af; font-size: 0.875em; border-top: 1px solid #f3f4f6; padding-top: 24px; }
    </style>
</head>
<body>
    <main>
        <h1>OAuth2 Provider</h1>
        <p><span class="status" role="status">System Operational</span></p>

        <div class="card">
            <h2>Available Endpoints</h2>
            <ul>
                <li>Authorize: <code>GET /authorize</code></li>
                <li>Token: <code>POST /token</code></li>
                <li>User Info: <code>GET /userinfo</code></li>
                <li>Register User: <code>POST /register</code></li>
                <li>Register Client: <code>POST /client/register</code></li>
            </ul>
        </div>

        <div class="footer">
            <p>Painted with 🎨 by Palette</p>
        </div>
    </main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
