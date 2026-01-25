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
	return c.HTML(http.StatusOK, `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Status</title><style>:root{--bg:#111;--text:#fff;--success:#4ade80}body{background:var(--bg);color:var(--text);font-family:system-ui,-apple-system,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}.status{text-align:center}.icon{width:64px;height:64px;color:var(--success)}@media(prefers-reduced-motion:reduce){*{animation:none!important}}</style></head><body><div class="status"><svg class="icon" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-label="System Operational"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg><h1>System Operational</h1></div></body></html>`)
}

func (h *HTMLHandler) LoginView(c echo.Context) error {
	continueTo := c.QueryParam("continue_to")
	return c.HTML(http.StatusOK, `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Login</title><style>:root{--bg:#111;--panel:#222;--text:#fff;--input:#333;--primary:#3b82f6}body{background:var(--bg);color:var(--text);font-family:system-ui,-apple-system,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}.card{background:var(--panel);padding:2rem;border-radius:8px;width:100%;max-width:320px;box-shadow:0 4px 6px rgba(0,0,0,0.3)}h2{margin-top:0}input{width:100%;padding:0.75rem;margin-bottom:1rem;background:var(--input);border:1px solid #444;color:var(--text);border-radius:4px;box-sizing:border-box}button{width:100%;padding:0.75rem;background:var(--primary);color:white;border:none;border-radius:4px;font-weight:600;cursor:pointer}button:hover{opacity:0.9}@media(prefers-reduced-motion:reduce){*{transition:none!important}}</style></head><body><div class="card"><h2>Login</h2><form action="/login" method="POST" onsubmit="return v()"><input type="hidden" name="continue_to" id="cont" value="`+continueTo+`"><input type="text" name="username" placeholder="Username" required aria-label="Username"><input type="password" name="password" placeholder="Password" required aria-label="Password"><button type="submit">Sign In</button></form><script>function v(){const u=document.getElementById('cont').value;if(u&&(u.startsWith('//')||!u.startsWith('/')))return false;return true}</script></div></body></html>`)
}
