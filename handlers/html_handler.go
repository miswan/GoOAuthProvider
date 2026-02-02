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
	html := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>System Status</title>
<style>
:root{--bg:#fff;--text:#111;--status:#16a34a}@media(prefers-color-scheme:dark){:root{--bg:#111;--text:#eee;--status:#4ade80}}
body{font-family:system-ui,-apple-system,sans-serif;background:var(--bg);color:var(--text);display:flex;align-items:center;justify-content:center;height:100vh;margin:0}
.container{text-align:center;padding:2rem}
h1{font-size:1.5rem;margin-bottom:0.5rem;display:flex;align-items:center;justify-content:center;gap:0.75rem}
.status-dot{width:12px;height:12px;background-color:var(--status);border-radius:50%}
@media(prefers-reduced-motion:reduce){*{transition:none!important}}
</style>
</head>
<body>
<main class="container">
<h1><span class="status-dot" aria-hidden="true"></span>System Operational</h1>
<p>The OAuth2 provider is running normally.</p>
</main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
