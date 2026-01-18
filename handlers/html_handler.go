package handlers

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
	"strings"
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
        :root { --bg: #111; --text: #eee; --success: #4ade80; }
        body { background: var(--bg); color: var(--text); font-family: system-ui, sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; }
        .status { text-align: center; }
        svg { width: 64px; height: 64px; color: var(--success); }
        @media (prefers-reduced-motion: reduce) { svg { animation: none; } }
    </style>
</head>
<body>
    <div class="status">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-label="System Operational">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h1>System Operational</h1>
    </div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}

func (h *HTMLHandler) Login(c echo.Context) error {
	next := c.QueryParam("next")

	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login</title>
    <style>
        :root { --bg: #111; --text: #eee; --input-bg: #222; --border: #333; --primary: #3b82f6; }
        body { background: var(--bg); color: var(--text); font-family: system-ui, sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; }
        form { background: var(--input-bg); padding: 2rem; border-radius: 0.5rem; border: 1px solid var(--border); width: 300px; }
        input { width: 100%; padding: 0.5rem; margin-bottom: 1rem; background: var(--bg); border: 1px solid var(--border); color: var(--text); box-sizing: border-box; }
        button { width: 100%; padding: 0.5rem; background: var(--primary); color: white; border: none; cursor: pointer; }
        @media (prefers-reduced-motion: reduce) { * { transition: none !important; animation: none !important; } }
    </style>
</head>
<body>
    <form action="/login" method="POST">
        <h2>Login</h2>
        <input type="hidden" name="next" value="{{.Next}}">
        <input type="text" name="username" placeholder="Username" required aria-label="Username">
        <input type="password" name="password" placeholder="Password" required aria-label="Password">
        <button type="submit">Sign In</button>
    </form>
</body>
</html>`

	t, err := template.New("login").Parse(tmpl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var buf strings.Builder
	if err := t.Execute(&buf, map[string]string{"Next": next}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.HTML(http.StatusOK, buf.String())
}
