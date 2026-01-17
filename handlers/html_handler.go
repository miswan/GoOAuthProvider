package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Index(c echo.Context) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>System Status</title>
<style>
body{background:#111;color:#eee;font-family:system-ui,-apple-system,sans-serif;display:grid;place-items:center;height:100vh;margin:0}
.status{text-align:center;animation:fade 0.5s ease-out}
@keyframes fade{from{opacity:0;transform:translateY(10px)}to{opacity:1;transform:translateY(0)}}
@media(prefers-reduced-motion:reduce){.status{animation:none}}
svg{width:64px;height:64px;color:#4ade80;margin-bottom:1rem}
</style>
</head>
<body>
<main class="status">
<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" aria-hidden="true">
<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
</svg>
<h1 aria-label="System Status: Operational">System Operational</h1>
</main>
</body></html>`
	return c.HTML(http.StatusOK, html)
}
