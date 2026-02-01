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
	// Minified HTML for system status page with accessible, dark-mode friendly design
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>OAuth2 Provider</title><style>:root{--bg:#111;--text:#eee;--success:#4ade80}@media(prefers-color-scheme:light){:root{--bg:#f9fafb;--text:#111827}}body{background:var(--bg);color:var(--text);font-family:system-ui,-apple-system,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}.status{display:flex;align-items:center;gap:.75rem;padding:1rem 1.5rem;border:1px solid #333;border-radius:99px}@media(prefers-color-scheme:light){.status{border-color:#e5e7eb}}.dot{width:.75rem;height:.75rem;background:var(--success);border-radius:50%;box-shadow:0 0 10px var(--success)}</style></head><body><div class="status" role="status" aria-label="System Status: Operational"><div class="dot" aria-hidden="true"></div><span>System Operational</span></div></body></html>`
	return c.HTML(http.StatusOK, html)
}
