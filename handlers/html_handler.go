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
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>OAuth2 Provider</title>
<style>:root{--primary:#4f46e5;--bg:#111;--text:#fff}@media(prefers-color-scheme:light){:root{--bg:#f9fafb;--text:#111827}}body{font-family:system-ui,-apple-system,sans-serif;background-color:var(--bg);color:var(--text);display:flex;align-items:center;justify-content:center;min-height:100vh;margin:0;line-height:1.5}.card{padding:2rem;text-align:center;border:1px solid rgba(255,255,255,0.1);border-radius:12px}@media(prefers-color-scheme:light){.card{border-color:rgba(0,0,0,0.1);background:white;box-shadow:0 4px 6px -1px rgba(0,0,0,0.1)}}.status{display:inline-flex;align-items:center;gap:0.5rem;color:#4ade80;font-weight:500;margin-bottom:1rem}h1{margin:0 0 0.5rem}p{color:#888;margin:0}</style></head>
<body><main class="card"><div class="status" role="status"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" stroke-linecap="round" stroke-linejoin="round"/><path d="M22 4L12 14.01l-3-3" stroke-linecap="round" stroke-linejoin="round"/></svg>System Operational</div><h1>OAuth2 Provider</h1><p>Ready to authenticate.</p></main></body></html>`
	return c.HTML(http.StatusOK, html)
}
