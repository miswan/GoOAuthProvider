package handlers

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Index(c echo.Context) error {
	return c.HTML(http.StatusOK, `
    <!DOCTYPE html>
    <html style="--bg:#111;background:var(--bg);color:#fff;font-family:sans-serif;">
    <head><title>OAuth2 Status</title></head>
    <body style="padding:2rem;">
        <h1>OAuth2 Provider Status</h1>
        <p>System is operational.</p>
        <p><a href="/login" style="color:#0af;">Go to Login</a></p>
    </body>
    </html>`)
}

func (h *HTMLHandler) Login(c echo.Context) error {
	next := c.QueryParam("next")
	html := `
    <!DOCTYPE html>
    <html style="--bg:#111;background:var(--bg);color:#fff;font-family:sans-serif;">
    <head><title>Login</title></head>
    <body style="display:flex;justify-content:center;align-items:center;height:100vh;flex-direction:column;">
        <h2>Login</h2>
        <form action="/login" method="POST" style="display:flex;flex-direction:column;gap:1em;width:300px;">
            <input type="text" name="username" placeholder="Username" required style="padding:0.5em;">
            <input type="password" name="password" placeholder="Password" required style="padding:0.5em;">
            <input type="hidden" name="next" value="{{.Next}}">
            <button type="submit" style="padding:0.5em;cursor:pointer;">Login</button>
        </form>
    </body>
    </html>
    `
	tmpl, err := template.New("login").Parse(html)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Template error")
	}

	data := map[string]interface{}{
		"Next": next,
	}

	return tmpl.Execute(c.Response().Writer, data)
}
