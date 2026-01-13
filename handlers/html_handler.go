package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler { return &HTMLHandler{} }

func (h *HTMLHandler) Index(c echo.Context) error {
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>OAuth2 Provider</title><style>:root{--bg:#111;--card:#1a1a1a;--text:#eee;--sub:#aaa;--border:#333;--acc:#4ade80}body{font-family:system-ui,-apple-system,sans-serif;background:var(--bg);color:var(--text);display:flex;align-items:center;justify-content:center;height:100vh;margin:0}.container{text-align:center;padding:2.5rem;border:1px solid var(--border);border-radius:12px;background:var(--card);box-shadow:0 4px 6px -1px rgba(0,0,0,0.1);max-width:400px;width:90%}h1{margin:0 0 .5rem;font-size:1.5rem;color:#fff}.status{color:var(--acc);font-weight:500;font-size:.875rem;display:inline-flex;align-items:center;gap:.5rem;background:rgba(74,222,128,.1);padding:.25rem .75rem;border-radius:99px;margin-bottom:1.5rem}.dot{width:8px;height:8px;background:currentColor;border-radius:50%;box-shadow:0 0 8px currentColor}p{color:var(--sub);margin:0;font-size:.95rem}.footer{margin-top:2rem;font-size:.75rem;color:#666;border-top:1px solid #333;padding-top:1rem}</style></head><body><div class="container"><div class="status" role="status" aria-live="polite"><span class="dot" aria-hidden="true"></span>System Operational</div><h1>OAuth2 Provider</h1><p>Secure identity and access management service is running.</p><div class="footer">v1.0.0</div></div></body></html>`
	return c.HTML(http.StatusOK, html)
}
