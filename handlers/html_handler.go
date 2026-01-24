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
<meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>OAuth2 Provider</title>
<style>
:root{--bg:#111;--text:#eee;--success:#4ade80}
body{font-family:system-ui,-apple-system,sans-serif;background:var(--bg);color:var(--text);display:grid;place-items:center;height:100vh;margin:0}
main{text-align:center;padding:2rem;border:1px solid #333;border-radius:12px;background:#1a1a1a}
.status{display:inline-flex;align-items:center;gap:0.5rem;margin-top:1rem;padding:0.5rem 1rem;background:#333;border-radius:99px;font-size:0.875rem}
.dot{width:8px;height:8px;background:var(--success);border-radius:50%;box-shadow:0 0 12px var(--success)}
</style>
</head>
<body>
<main role="main">
<h1>OAuth2 Provider</h1>
<div class="status" role="status" aria-label="System status: Operational">
<span class="dot"></span>Operational
</div>
</main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
