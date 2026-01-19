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
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>System Status</title><style>:root{--bg:#111;--text:#fff;--success:#4ade80}body{background:var(--bg);color:var(--text);font-family:system-ui,-apple-system,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}@media(prefers-reduced-motion:reduce){*{animation:none!important;transition:none!important}}</style></head><body><div style="display:flex;gap:0.5rem;align-items:center"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="var(--success)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg><span>System Operational</span></div></body></html>`
	return c.HTML(http.StatusOK, html)
}

func (h *HTMLHandler) Login(c echo.Context) error {
	next := c.QueryParam("next")
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Login</title><style>:root{--bg:#111;--text:#fff}body{background:var(--bg);color:var(--text);font-family:system-ui,-apple-system,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}form{display:flex;flex-direction:column;gap:1rem}input{padding:0.5rem}button{padding:0.5rem;cursor:pointer}@media(prefers-reduced-motion:reduce){*{animation:none!important;transition:none!important}}</style></head><body><form method="POST" action="/login"><input type="hidden" name="next" value="` + next + `"><input type="text" name="username" placeholder="Username" required><input type="password" name="password" placeholder="Password" required><button type="submit">Login</button></form></body></html>`
	return c.HTML(http.StatusOK, html)
}
