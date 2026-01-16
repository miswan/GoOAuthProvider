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
	// Simple, dark-themed status page with a pulse animation
	html := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>OAuth2 Provider</title>
<style>
:root{--bg:#111;--txt:#eee;--acc:#3b82f6}body{background:var(--bg);color:var(--txt);font-family:system-ui,sans-serif;display:grid;place-items:center;height:100vh;margin:0}
.s{display:flex;align-items:center;gap:.5rem}.d{width:12px;height:12px;background:#10b981;border-radius:50%;box-shadow:0 0 12px #10b981;animation:p 2s infinite}
@keyframes p{0%{opacity:1}50%{opacity:.5}100%{opacity:1}}
@media(prefers-reduced-motion:reduce){.d{animation:none}}
</style>
</head>
<body>
<main>
<div class="s"><div class="d"></div><h1>System Operational</h1></div>
<p>OAuth2 Provider is running smoothly.</p>
</main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
