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
<meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>OAuth2 Provider</title>
<style>
:root { --bg: #111; --text: #eee; --accent: #3b82f6; --success: #4ade80; }
@media (prefers-color-scheme: light) { :root { --bg: #fff; --text: #111; --accent: #2563eb; --success: #16a34a; } }
body { font-family: system-ui, -apple-system, sans-serif; background: var(--bg); color: var(--text); display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; }
.card { padding: 2rem; border: 1px solid #333; border-radius: 12px; max-width: 400px; width: 100%; }
@media (prefers-color-scheme: light) { .card { border-color: #eee; box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1); } }
h1 { margin: 0 0 1rem; font-size: 1.5rem; display: flex; align-items: center; gap: 0.5rem; }
.status { width: 10px; height: 10px; background: var(--success); border-radius: 50%; display: inline-block; }
p { color: #888; margin-bottom: 1.5rem; line-height: 1.5; }
.links { display: flex; flex-direction: column; gap: 0.5rem; }
a { color: var(--accent); text-decoration: none; font-weight: 500; }
a:hover { text-decoration: underline; }
@media (prefers-reduced-motion: reduce) { *, ::before, ::after { animation-duration: 0.01ms !important; animation-iteration-count: 1 !important; transition-duration: 0.01ms !important; scroll-behavior: auto !important; } }
</style>
</head>
<body>
<div class="card">
<h1><span class="status" aria-hidden="true"></span> System Operational</h1>
<p>OAuth2 Provider is running. Use the endpoints below to interact with the service.</p>
<div class="links">
<a href="/authorize">Authorize Endpoint</a>
<a href="/userinfo">User Info</a>
</div>
</div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
