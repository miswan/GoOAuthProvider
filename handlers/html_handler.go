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

func (h *HTMLHandler) LoginView(c echo.Context) error {
	redirectTo := c.QueryParam("redirect_to")

	const tpl = `
<!DOCTYPE html>
<html>
<head>
    <title>Login - OAuth Provider</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; }
        .container { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); width: 300px; }
        h2 { text-align: center; color: #1a73e8; }
        form { display: flex; flex-direction: column; gap: 1rem; }
        input { padding: 0.5rem; border: 1px solid #ddd; border-radius: 4px; }
        button { padding: 0.5rem; background: #1a73e8; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: #1557b0; }
        .error { color: red; font-size: 0.9rem; text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Login</h2>
        <form action="/login" method="POST">
            <input type="hidden" name="redirect_to" value="{{.RedirectTo}}">
            <input type="text" name="username" placeholder="Username" required>
            <input type="password" name="password" placeholder="Password" required>
            <button type="submit">Sign In</button>
        </form>
    </div>
</body>
</html>`

	t, err := template.New("login").Parse(tpl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	data := struct {
		RedirectTo string
	}{
		RedirectTo: redirectTo,
	}

	return t.Execute(c.Response().Writer, data)
}

func (h *HTMLHandler) Index(c echo.Context) error {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>OAuth Provider</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; }
        .container { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); text-align: center; }
        h1 { color: #1a73e8; }
        p { color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <h1>OAuth2 Provider</h1>
        <p>Service is running.</p>
    </div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
