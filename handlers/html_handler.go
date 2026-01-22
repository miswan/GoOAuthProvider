package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Index(c echo.Context) error {
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>System Operational</title><style>:root{--bg:#111;--text:#eee;--success:#4ade80}body{font-family:system-ui,-apple-system,sans-serif;background:var(--bg);color:var(--text);display:flex;align-items:center;justify-content:center;height:100vh;margin:0}.status{text-align:center}svg{width:48px;height:48px;color:var(--success)}@media(prefers-reduced-motion:reduce){*{animation:none!important;transition:none!important}}</style></head><body><div class="status"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg><h1>System Operational</h1></div></body></html>`
	return c.HTML(http.StatusOK, html)
}

func (h *HTMLHandler) Login(c echo.Context) error {
	returnTo := c.QueryParam("return_to")
	if returnTo == "" {
		returnTo = "/"
	}

	// Simple escaping for JS context, not robust but okay for this scope
	returnToEscaped := url.QueryEscape(returnTo)

	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Login</title><style>:root{--bg:#111;--text:#eee;--input-bg:#222;--border:#333;--primary:#3b82f6}body{font-family:system-ui,-apple-system,sans-serif;background:var(--bg);color:var(--text);display:flex;align-items:center;justify-content:center;height:100vh;margin:0}form{display:flex;flex-direction:column;gap:1rem;width:300px}input{padding:0.5rem;background:var(--input-bg);border:1px solid var(--border);color:var(--text);border-radius:4px}button{padding:0.5rem;background:var(--primary);color:white;border:none;border-radius:4px;cursor:pointer}a{color:var(--primary);text-align:center;font-size:0.9rem}@media(prefers-reduced-motion:reduce){*{animation:none!important;transition:none!important}}</style></head><body><form id="loginForm"><h1>Login</h1><input type="text" name="username" placeholder="Username" required><input type="password" name="password" placeholder="Password" required><button type="submit">Login</button><a href="/register">Register</a></form><script>const returnTo = decodeURIComponent("` + returnToEscaped + `");document.getElementById('loginForm').onsubmit=async(e)=>{e.preventDefault();const formData=new FormData(e.target);const res=await fetch('/login',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(Object.fromEntries(formData))});if(res.ok){window.location.href=returnTo;}else{alert('Login failed');}}</script></body></html>`
	return c.HTML(http.StatusOK, html)
}

func (h *HTMLHandler) Register(c echo.Context) error {
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Register</title><style>:root{--bg:#111;--text:#eee;--input-bg:#222;--border:#333;--primary:#3b82f6}body{font-family:system-ui,-apple-system,sans-serif;background:var(--bg);color:var(--text);display:flex;align-items:center;justify-content:center;height:100vh;margin:0}form{display:flex;flex-direction:column;gap:1rem;width:300px}input{padding:0.5rem;background:var(--input-bg);border:1px solid var(--border);color:var(--text);border-radius:4px}button{padding:0.5rem;background:var(--primary);color:white;border:none;border-radius:4px;cursor:pointer}a{color:var(--primary);text-align:center;font-size:0.9rem}@media(prefers-reduced-motion:reduce){*{animation:none!important;transition:none!important}}</style></head><body><form id="registerForm"><h1>Register</h1><input type="text" name="username" placeholder="Username" required><input type="email" name="email" placeholder="Email" required><input type="password" name="password" placeholder="Password" required><button type="submit">Register</button><a href="/login">Login</a></form><script>document.getElementById('registerForm').onsubmit=async(e)=>{e.preventDefault();const formData=new FormData(e.target);const res=await fetch('/register',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(Object.fromEntries(formData))});if(res.ok){window.location.href='/login';}else{alert('Registration failed');}}</script></body></html>`
	return c.HTML(http.StatusOK, html)
}
