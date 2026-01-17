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
	return c.HTML(http.StatusOK, `
<!DOCTYPE html>
<html>
<head>
    <title>OAuth2 Provider</title>
    <style>body{font-family:sans-serif;padding:2rem;text-align:center;background:#111;color:#eee} a{color:#4daafc}</style>
</head>
<body>
    <h1>System Operational</h1>
    <p><a href="/login">Login</a> | <a href="/register">Register</a></p>
</body>
</html>`)
}

func (h *HTMLHandler) Login(c echo.Context) error {
	next := c.QueryParam("next")
	if next == "" {
		next = "/"
	}

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background: #111; color: #eee; margin: 0; }
        .container { background: #222; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.3); width: 300px; }
        h2 { margin-top: 0; text-align: center; }
        input { display: block; margin-bottom: 1rem; width: 100%; padding: 0.5rem; box-sizing: border-box; background: #333; border: 1px solid #444; color: white; border-radius: 4px; }
        button { background: #007bff; color: white; border: none; padding: 0.7rem; width: 100%; cursor: pointer; border-radius: 4px; font-weight: bold; }
        button:hover { background: #0056b3; }
        .link { text-align: center; margin-top: 1rem; font-size: 0.9rem; }
        a { color: #4daafc; text-decoration: none; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Login</h2>
        <form action="/login" method="POST">
            <input type="hidden" name="next" value="` + next + `">
            <input type="text" name="username" placeholder="Username" required>
            <input type="password" name="password" placeholder="Password" required>
            <button type="submit">Sign In</button>
        </form>
        <div class="link"><a href="/register">Create Account</a></div>
    </div>
</body>
</html>
`
	return c.HTML(http.StatusOK, html)
}

func (h *HTMLHandler) Register(c echo.Context) error {
    html := `
<!DOCTYPE html>
<html>
<head>
    <title>Register</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background: #111; color: #eee; margin: 0; }
        .container { background: #222; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.3); width: 300px; }
        h2 { margin-top: 0; text-align: center; }
        input { display: block; margin-bottom: 1rem; width: 100%; padding: 0.5rem; box-sizing: border-box; background: #333; border: 1px solid #444; color: white; border-radius: 4px; }
        button { background: #28a745; color: white; border: none; padding: 0.7rem; width: 100%; cursor: pointer; border-radius: 4px; font-weight: bold; }
        button:hover { background: #218838; }
        .link { text-align: center; margin-top: 1rem; font-size: 0.9rem; }
        a { color: #4daafc; text-decoration: none; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Register</h2>
        <form action="/register" method="POST">
            <input type="text" name="username" placeholder="Username" required>
            <input type="email" name="email" placeholder="Email (optional)">
            <input type="password" name="password" placeholder="Password" required>
            <button type="submit">Register</button>
        </form>
        <div class="link"><a href="/login">Back to Login</a></div>
    </div>
</body>
</html>
`
    return c.HTML(http.StatusOK, html)
}
