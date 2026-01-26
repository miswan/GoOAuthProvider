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
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>System Status</title><style>:root{--bg:#111;--text:#eee;--success:#4ade80}@media(prefers-color-scheme:light){:root{--bg:#fff;--text:#111}}body{background:var(--bg);color:var(--text);font-family:system-ui,-apple-system,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}.status{display:flex;align-items:center;gap:.5rem;font-weight:500}svg{width:1.5rem;height:1.5rem;color:var(--success)}</style></head><body><div class="status" role="status" aria-label="System status: Operational"><svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg><span>System Operational</span></div></body></html>`
	return c.HTML(http.StatusOK, html)
}
