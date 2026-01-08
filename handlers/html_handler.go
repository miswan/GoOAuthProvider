package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler { return &HTMLHandler{} }

func (h *HTMLHandler) Welcome(c echo.Context) error {
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>OAuth2 Provider - System Status</title>
<style>
:root { --bg: #1a1a1a; --text: #f0f0f0; --card: #2a2a2a; }
@media(prefers-color-scheme:light) { :root { --bg: #f9fafb; --text: #111827; --card: #fff; } }
body { font-family: system-ui, sans-serif; background: var(--bg); color: var(--text); display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
.card { background: var(--card); padding: 2rem; border-radius: .5rem; box-shadow: 0 4px 6px -1px rgba(0,0,0,.1); text-align: center; width: 100%; max-width: 24rem; }
.status { display: inline-flex; align-items: center; background: rgba(16,185,129,.1); color: #10b981; padding: .25rem .75rem; border-radius: 99px; font-weight: 500; margin-bottom: 1rem; }
.dot { width: .5rem; height: .5rem; background: currentColor; border-radius: 50%; margin-right: .5rem; }
p { color: #6b7280; } @media(prefers-color-scheme:dark) { p { color: #9ca3af; } }
</style></head>
<body>
<div class="card">
    <h1 style="margin-top:0;font-size:1.5rem">OAuth2 Provider</h1>
    <div class="status" role="status"><span class="dot"></span>System Operational</div>
    <p>The OAuth2 provider service is running and ready to accept connections.</p>
</div></body></html>`
	return c.HTML(http.StatusOK, html)
}
