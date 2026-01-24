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
	return c.HTML(http.StatusOK, `<!DOCTYPE html><html style="background:#111;color:#fff;font-family:system-ui,-apple-system,sans-serif"><head><title>System Operational</title><style>@media(prefers-reduced-motion:reduce){*{animation:none!important}}</style></head><body style="display:flex;justify-content:center;align-items:center;height:100vh;margin:0"><div style="text-align:center"><svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="#4ade80" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-label="Success"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg><h1>System Operational</h1></div></body></html>`)
}

func (h *HTMLHandler) Login(c echo.Context) error {
	return c.HTML(http.StatusOK, `<!DOCTYPE html><html style="background:#111;color:#fff;font-family:system-ui,-apple-system,sans-serif"><head><title>Login</title><style>input{display:block;margin:10px 0;padding:8px;width:100%;box-sizing:border-box}button{padding:10px;width:100%;cursor:pointer;background:#4ade80;border:none;font-weight:bold}@media(prefers-reduced-motion:reduce){*{animation:none!important}}</style></head><body style="display:flex;justify-content:center;align-items:center;height:100vh;margin:0"><form style="width:300px" onsubmit="event.preventDefault();fetch('/login',{method:'POST',headers:{'Content-Type':'application/x-www-form-urlencoded'},body:new URLSearchParams(new FormData(this))}).then(r=>{if(r.ok){let u=new URLSearchParams(location.search).get('continue_to');if(u&&u.startsWith('/')&&!u.startsWith('//'))location.href=u;else location.href='/'}else{alert('Login failed')}})"><h2>Login</h2><input name="username" placeholder="Username" required><input type="password" name="password" placeholder="Password" required><button>Login</button><p><a href="/register" style="color:#4ade80">Register</a></p></form></body></html>`)
}

func (h *HTMLHandler) Register(c echo.Context) error {
	return c.HTML(http.StatusOK, `<!DOCTYPE html><html style="background:#111;color:#fff;font-family:system-ui,-apple-system,sans-serif"><head><title>Register</title><style>input{display:block;margin:10px 0;padding:8px;width:100%;box-sizing:border-box}button{padding:10px;width:100%;cursor:pointer;background:#4ade80;border:none;font-weight:bold}@media(prefers-reduced-motion:reduce){*{animation:none!important}}</style></head><body style="display:flex;justify-content:center;align-items:center;height:100vh;margin:0"><form style="width:300px" onsubmit="event.preventDefault();fetch('/register',{method:'POST',headers:{'Content-Type':'application/x-www-form-urlencoded'},body:new URLSearchParams(new FormData(this))}).then(r=>{if(r.ok){alert('Registered!');location.href='/login'}else{alert('Registration failed')}})"><h2>Register</h2><input name="username" placeholder="Username" required><input type="email" name="email" placeholder="Email" required><input type="password" name="password" placeholder="Password" required><button>Register</button><p><a href="/login" style="color:#4ade80">Login</a></p></form></body></html>`)
}
