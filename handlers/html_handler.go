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
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>System Status</title><style>:root{--bg:#111;--text:#fff;--success:#4ade80;--font:system-ui,-apple-system,sans-serif}body{background:var(--bg);color:var(--text);font-family:var(--font);display:flex;align-items:center;justify-content:center;height:100vh;margin:0}.status{display:flex;align-items:center;gap:0.5rem}.dot{width:10px;height:10px;background:var(--success);border-radius:50%}@media(prefers-reduced-motion:reduce){*{transition:none!important;animation:none!important}}</style></head><body><div class="status"><div class="dot"></div>System Operational</div></body></html>`
	return c.HTML(http.StatusOK, html)
}

func (h *HTMLHandler) Login(c echo.Context) error {
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Login</title><style>:root{--bg:#111;--text:#fff;--input-bg:#222;--border:#333;--font:system-ui,-apple-system,sans-serif;--primary:#3b82f6}body{background:var(--bg);color:var(--text);font-family:var(--font);display:flex;align-items:center;justify-content:center;height:100vh;margin:0}form{display:flex;flex-direction:column;gap:1rem;width:300px;padding:2rem;border:1px solid var(--border);border-radius:8px}input{background:var(--input-bg);border:1px solid var(--border);color:var(--text);padding:0.5rem;border-radius:4px}button{background:var(--primary);color:#fff;border:none;padding:0.5rem;border-radius:4px;cursor:pointer}@media(prefers-reduced-motion:reduce){*{transition:none!important;animation:none!important}}</style></head><body><form id="loginForm"><h2>Login</h2><input type="text" name="username" placeholder="Username" required><input type="password" name="password" placeholder="Password" required><button type="submit">Login</button><div id="error" style="color:red;display:none"></div></form><script>document.getElementById('loginForm').onsubmit=async(e)=>{e.preventDefault();const t=new FormData(e.target),n=Object.fromEntries(t.entries()),r=document.getElementById('error');r.style.display='none';try{const o=await fetch('/login',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(n)});if(o.ok){const s=new URLSearchParams(window.location.search).get('continue_to')||'/';window.location.href=(s.startsWith('/')&&!s.startsWith('//'))?s:'/'}else{const a=await o.json();r.innerText=a.message||'Login failed',r.style.display='block'}}catch(i){r.innerText='An error occurred',r.style.display='block'}};</script></body></html>`
	return c.HTML(http.StatusOK, html)
}
