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

func (h *HTMLHandler) Login(c echo.Context) error {
	next := c.QueryParam("next")

	// Safe template to prevent XSS
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Login - OAuth2 Provider</title>
    <style>
        :root { --bg: #111; --text: #eee; --input-bg: #222; --border: #333; }
        body { background: var(--bg); color: var(--text); font-family: system-ui, sans-serif; display: grid; place-items: center; height: 100vh; margin: 0; }
        form { background: var(--input-bg); padding: 2rem; border-radius: 8px; border: 1px solid var(--border); width: 300px; }
        input { display: block; width: 100%; margin-bottom: 1rem; padding: 0.5rem; background: var(--bg); border: 1px solid var(--border); color: var(--text); box-sizing: border-box; }
        button { width: 100%; padding: 0.5rem; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: #0056b3; }
    </style>
</head>
<body>
    <form action="/login" method="POST">
		<input type="hidden" name="next" value="{{.Next}}">
        <h2 style="margin-top:0">Login</h2>
        <input type="text" name="username" placeholder="Username" required>
        <input type="password" name="password" placeholder="Password" required>
        <button type="submit">Login</button>
    </form>
</body>
</html>`

	t, err := template.New("login").Parse(tmpl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Template error")
	}

	data := map[string]string{
		"Next": next,
	}

	return t.Execute(c.Response().Writer, data)
}

func (h *HTMLHandler) Index(c echo.Context) error {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>OAuth2 Provider Status</title>
    <style>
        body { background: #111; color: #eee; font-family: system-ui, sans-serif; padding: 2rem; }
        .status { padding: 1rem; background: #222; border-radius: 4px; display: inline-block; }
        .ok { color: #4cd964; }
    </style>
</head>
<body>
    <h1>System Status</h1>
    <div class="status">
        Operational <span class="ok">●</span>
    </div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
