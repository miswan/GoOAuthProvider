package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Welcome(c echo.Context) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>OAuth2 Provider Status</title>
<style>
:root{--primary:#2563eb;--bg:#f8fafc;--text:#1e293b;--card:#fff;--ok:#16a34a}
body{font-family:system-ui,-apple-system,sans-serif;background:var(--bg);color:var(--text);line-height:1.5;margin:0;padding:2rem;display:grid;place-items:center;min-height:100vh}
main{background:var(--card);padding:2rem;border-radius:12px;box-shadow:0 4px 6px -1px rgb(0 0 0 / 0.1);max-width:32rem;width:100%}
h1{margin:0 0 0.5rem;color:var(--primary);font-size:1.5rem}
.badge{display:inline-flex;align-items:center;background:#dcfce7;color:var(--ok);padding:.25rem .75rem;border-radius:99px;font-size:.875rem;font-weight:600;margin-bottom:1.5rem}
.badge::before{content:"";width:6px;height:6px;background:currentColor;border-radius:50%;margin-right:.5rem}
ul{list-style:none;padding:0;margin:1.5rem 0}
li{padding:.5rem 0;border-bottom:1px solid #f1f5f9;display:flex;justify-content:space-between;font-size:0.875rem}
code{background:#f1f5f9;padding:.2rem .4rem;border-radius:4px;color:#475569;font-family:ui-monospace,monospace}
footer{margin-top:2rem;font-size:.75rem;color:#94a3b8;text-align:center}
</style>
</head>
<body>
<main>
<div class="badge" role="status">System Operational</div>
<h1>OAuth2 Provider</h1>
<p>Ready to handle authentication requests.</p>
<nav aria-label="Available Endpoints">
<ul>
<li><span>Authorize</span><code>GET /authorize</code></li>
<li><span>Token</span><code>POST /token</code></li>
<li><span>Register</span><code>POST /register</code></li>
<li><span>Login</span><code>POST /login</code></li>
</ul>
</nav>
<footer>
Palette UX • OAuth2 Provider
</footer>
</main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
