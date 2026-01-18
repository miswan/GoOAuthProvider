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
	const html = `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>System Status</title><style>:root{--bg:#111;--text:#fff;--success:#4ade80}body{background:var(--bg);color:var(--text);font-family:system-ui,-apple-system,sans-serif;display:grid;place-items:center;height:100vh;margin:0}@media(prefers-reduced-motion:reduce){*{animation:none!important;transition:none!important}}</style></head><body><div style="text-align:center"><svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="var(--success)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" role="img" aria-label="Status: Operational"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg><h1 style="margin-top:1rem;font-size:1.5rem">System Operational</h1></div></body></html>`
	return c.HTML(http.StatusOK, html)
}
