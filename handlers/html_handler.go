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
        :root { --bg: #111; --text: #eee; --accent: #3b82f6; }
        body { background: var(--bg); color: var(--text); font-family: system-ui, sans-serif; display: grid; place-items: center; height: 100vh; margin: 0; }
        main { text-align: center; padding: 2rem; border: 1px solid #333; border-radius: 8px; }
        h1 { margin-bottom: 0.5rem; }
        .status { color: #22c55e; font-weight: bold; display: inline-flex; align-items: center; gap: 0.5rem; }
        .status::before { content: ""; width: 8px; height: 8px; background: currentColor; border-radius: 50%; display: inline-block; }
    </style>
</head>
<body>
    <main>
        <h1>OAuth2 Provider</h1>
        <p class="status" role="status">System Operational</p>
        <p>Ready to handle authentication requests.</p>
    </main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
