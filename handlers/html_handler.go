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
		:root { --bg: #111; --text: #eee; --accent: #4ade80; }
		@media (prefers-color-scheme: light) { :root { --bg: #f9fafb; --text: #111827; } }
		body { font-family: system-ui, -apple-system, sans-serif; background: var(--bg); color: var(--text); display: grid; place-items: center; height: 100vh; margin: 0; }
		.status { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1.5rem; background: rgba(128,128,128,0.1); border-radius: 99px; border: 1px solid rgba(128,128,128,0.2); font-weight: 500; }
		.dot { width: 10px; height: 10px; background: var(--accent); border-radius: 50%; box-shadow: 0 0 12px var(--accent); animation: pulse 2s infinite; }
		@media (prefers-reduced-motion: reduce) { .dot { animation: none; } }
		@keyframes pulse { 0% { opacity: 1; } 50% { opacity: 0.5; } 100% { opacity: 1; } }
	</style>
</head>
<body>
	<div class="status" role="status" aria-live="polite">
		<div class="dot" aria-hidden="true"></div>
		<span>System Operational</span>
	</div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
