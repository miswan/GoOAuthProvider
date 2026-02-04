package handlers

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Index(c echo.Context) error {
	tmpl := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>OAuth2 Provider</title><style>:root{--bg:#111;--text:#eee;--success:#4ade80}@media(prefers-color-scheme:light){:root{--bg:#fff;--text:#111;--success:#16a34a}}body{font-family:system-ui,-apple-system,sans-serif;background:var(--bg);color:var(--text);display:flex;align-items:center;justify-content:center;height:100vh;margin:0}.status{display:flex;align-items:center;gap:0.5rem;font-size:1.25rem}.dot{width:12px;height:12px;background:var(--success);border-radius:50%;box-shadow:0 0 8px var(--success)}@media(prefers-reduced-motion:reduce){*{animation:none!important;transition:none!important}}</style></head><body><main class="status" role="status" aria-live="polite"><div class="dot" aria-hidden="true"></div><span>System Operational</span></main></body></html>`

	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return t.Execute(c.Response().Writer, nil)
}
