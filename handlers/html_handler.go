package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler { return &HTMLHandler{} }

func (h *HTMLHandler) Welcome(c echo.Context) error {
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>OAuth2 Provider Service</title><style>
:root{--p:#2563eb;--t:#1f2937;--b:#f3f4f6;--c:#fff}body{font-family:system-ui,-apple-system,sans-serif;line-height:1.5;color:var(--t);background:var(--b);margin:0;display:flex;justify-content:center;align-items:center;min-height:100vh;padding:1rem}main{background:var(--c);padding:2rem;border-radius:.5rem;box-shadow:0 4px 6px -1px #0000001a;max-width:32rem;width:100%}h1{font-size:1.5rem;font-weight:700;margin:0 0 1rem;color:var(--p)}p{margin-bottom:1.5rem}.status{display:inline-flex;align-items:center;padding:.25rem .75rem;background:#d1fae5;color:#065f46;border-radius:99px;font-size:.875rem;font-weight:500}.status::before{content:"";display:block;width:.5rem;height:.5rem;background:#059669;border-radius:50%;margin-right:.5rem}footer{margin-top:2rem;font-size:.875rem;color:#6b7280;text-align:center}</style></head>
<body><main><h1>OAuth2 Provider</h1><p>Welcome to the OAuth2 Provider service. This API handles user authentication and authorization using the OAuth 2.0 protocol.</p><div class="status">System Operational</div><footer><p>Ready to handle requests.</p></footer></main></body></html>`
	return c.HTML(http.StatusOK, html)
}
