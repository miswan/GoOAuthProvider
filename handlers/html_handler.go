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
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>System Status</title><style>:root{--bg:#111;--text:#eee;--success:#4ade80}body{background:var(--bg);color:var(--text);font-family:system-ui,-apple-system,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}.status{display:flex;align-items:center;gap:0.5rem;font-weight:500}.dot{width:8px;height:8px;background:var(--success);border-radius:50%;box-shadow:0 0 8px var(--success)}</style></head><body><div class="status" role="status" aria-label="System status: Operational"><div class="dot"></div>System Operational</div></body></html>`
	return c.HTML(http.StatusOK, html)
}
