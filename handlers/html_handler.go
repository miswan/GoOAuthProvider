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
:root { --bg: #111; --fg: #eee; --accent: #4ade80; }
@media (prefers-color-scheme: light) { :root { --bg: #fff; --fg: #111; --accent: #16a34a; } }
body { font-family: system-ui, -apple-system, sans-serif; background: var(--bg); color: var(--fg); display: grid; place-items: center; height: 100vh; margin: 0; }
.status { display: flex; align-items: center; gap: 0.5rem; font-weight: 500; }
.dot { width: 0.75rem; height: 0.75rem; background: var(--accent); border-radius: 50%; }
</style>
</head>
<body>
<div class="status" role="status" aria-live="polite">
<div class="dot"></div>
<span>System Operational</span>
</div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
