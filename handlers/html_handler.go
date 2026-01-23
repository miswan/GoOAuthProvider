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
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>System Operational</title><style>:root{--bg:#111;--text:#eee;--success:#4ade80}body{background:var(--bg);color:var(--text);font-family:system-ui,-apple-system,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}.status{text-align:center}.icon{color:var(--success);width:48px;height:48px}h1{margin:1rem 0;font-size:1.5rem;font-weight:500}p{opacity:.8}</style></head><body><div class="status"><svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-label="System operational"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg><h1>System Operational</h1><p>OAuth2 Provider is running.</p></div></body></html>`
	return c.HTML(http.StatusOK, html)
}
