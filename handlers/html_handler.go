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
	// Minified HTML string with dark mode and success icon
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>System Operational</title><style>:root{--bg:#111;--text:#eee;--success:#4ade80}body{background-color:var(--bg);color:var(--text);font-family:system-ui,-apple-system,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}@media(prefers-reduced-motion:reduce){*{animation:none!important;transition:none!important}}svg{width:48px;height:48px;color:var(--success)}</style></head><body><div style="text-align:center"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg><h1>System Operational</h1></div></body></html>`
	return c.HTML(http.StatusOK, html)
}
