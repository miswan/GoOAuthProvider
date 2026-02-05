package handlers

import (
	"github.com/labstack/echo/v4"
	"html/template"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Index(c echo.Context) error {
	tmpl, err := template.New("index").Parse(indexHTML)
	if err != nil {
		return err
	}
	return tmpl.Execute(c.Response().Writer, nil)
}

func (h *HTMLHandler) Login(c echo.Context) error {
	data := map[string]interface{}{
		"return_to": c.QueryParam("return_to"),
	}
	tmpl, err := template.New("login").Parse(loginHTML)
	if err != nil {
		return err
	}
	return tmpl.Execute(c.Response().Writer, data)
}

const indexHTML = `<!DOCTYPE html>
<html>
<head><title>OAuth2 Provider</title></head>
<body>
<h1>OAuth2 Provider</h1>
<p>System operational.</p>
<p><a href="/login">Login</a></p>
</body>
</html>`

const loginHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; max-width: 400px; margin: 40px auto; padding: 20px; line-height: 1.5; }
        .form-group { margin-bottom: 15px; }
        label { display: block; margin-bottom: 5px; font-weight: 500; }
        input { width: 100%; padding: 8px; box-sizing: border-box; border: 1px solid #ccc; border-radius: 4px; }
        button { width: 100%; padding: 10px; background: #2563eb; color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: 500; }
        button:hover { background: #1d4ed8; }
    </style>
</head>
<body>
<h1>Login</h1>
<form action="/login" method="POST">
    {{if .return_to}}
    <input type="hidden" name="return_to" value="{{.return_to}}">
    {{end}}
    <div class="form-group">
        <label for="username">Username</label>
        <input type="text" id="username" name="username" required>
    </div>
    <div class="form-group">
        <label for="password">Password</label>
        <input type="password" id="password" name="password" required>
    </div>
    <button type="submit">Sign In</button>
</form>
</body>
</html>`
