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
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>System Status</title><style>:root{--bg:#111;--text:#fff;--success:#4ade80}body{background:var(--bg);color:var(--text);font-family:system-ui,-apple-system,sans-serif;display:flex;justify-content:center;align-items:center;height:100vh;margin:0}h1{display:flex;align-items:center;gap:10px}svg{width:24px;height:24px;color:var(--success)}@media(prefers-reduced-motion:reduce){*{animation:none!important;transition:none!important}}</style></head><body><h1><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>System Operational</h1></body></html>`
	return c.HTML(http.StatusOK, html)
}

func (h *HTMLHandler) Login(c echo.Context) error {
	redirect := c.QueryParam("redirect_to")
	if redirect == "" {
		redirect = "/"
	}
	// Simple login form with embedded CSS and accessibility features
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Login</title><style>:root{--bg:#111;--text:#fff;--primary:#3b82f6}body{background:var(--bg);color:var(--text);font-family:system-ui,-apple-system,sans-serif;display:flex;justify-content:center;align-items:center;height:100vh;margin:0}form{display:flex;flex-direction:column;gap:1rem;padding:2rem;border:1px solid #333;border-radius:8px;width:300px}input{padding:0.5rem;background:#222;border:1px solid #444;color:#fff;border-radius:4px}button{padding:0.5rem;background:var(--primary);color:#fff;border:none;border-radius:4px;cursor:pointer}@media(prefers-reduced-motion:reduce){*{transition:none!important}}</style></head><body><form id="loginForm" method="POST" action="/login?redirect_to=`+redirect+`"><h2>Login</h2><input type="text" name="username" placeholder="Username" required aria-label="Username"><input type="password" name="password" placeholder="Password" required aria-label="Password"><button type="submit">Login</button></form><script>document.getElementById('loginForm').addEventListener('submit',async e=>{e.preventDefault();const f=e.target;const r=await fetch(f.action,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(Object.fromEntries(new FormData(f)))});if(r.ok){const d=await r.json();localStorage.setItem('token',d.token);window.location.href=new URLSearchParams(window.location.search).get('redirect_to')||'/';}else{alert('Login failed')}});</script></body></html>`
	return c.HTML(http.StatusOK, html)
}
