package handlers

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Index(c echo.Context) error {
	return c.HTML(http.StatusOK, "<h1>OAuth2 Provider</h1><p>System Operational</p>")
}

func (h *HTMLHandler) Login(c echo.Context) error {
	returnTo := c.QueryParam("return_to")

	htmlTmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background: #f0f2f5; }
        .login-container { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1); width: 100%; max-width: 400px; }
        h2 { text-align: center; color: #1f2937; }
        .form-group { margin-bottom: 1rem; }
        label { display: block; margin-bottom: 0.5rem; color: #4b5563; }
        input { width: 100%; padding: 0.5rem; border: 1px solid #d1d5db; border-radius: 4px; box-sizing: border-box; }
        button { width: 100%; padding: 0.75rem; background: #2563eb; color: white; border: none; border-radius: 4px; font-weight: bold; cursor: pointer; }
        button:hover { background: #1d4ed8; }
    </style>
</head>
<body>
    <div class="login-container">
        <h2>Login</h2>
        <form action="/login" method="POST">
            <input type="hidden" name="return_to" value="{{.ReturnTo}}">
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
    </div>
</body>
</html>
`
	t, err := template.New("login").Parse(htmlTmpl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]string{"ReturnTo": returnTo}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.HTML(http.StatusOK, buf.String())
}
