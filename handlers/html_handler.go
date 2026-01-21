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
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>System Status</title>
    <style>
        :root { --bg: #111; --text: #eee; --success: #4ade80; }
        body { background: var(--bg); color: var(--text); font-family: system-ui, -apple-system, sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; }
        .status { text-align: center; }
        .dot { height: 12px; width: 12px; background-color: var(--success); border-radius: 50%; display: inline-block; margin-right: 8px; }
    </style>
</head>
<body>
    <div class="status">
        <h1><span class="dot" aria-hidden="true"></span>System Operational</h1>
        <p>OAuth2 Provider is running.</p>
    </div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
