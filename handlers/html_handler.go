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
	return c.HTML(http.StatusOK, "<h1>System Operational</h1><p><a href='/login'>Login</a></p>")
}

func (h *HTMLHandler) Login(c echo.Context) error {
	returnTo := c.QueryParam("return_to")
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #f0f0f0; }
        .login-box { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); width: 300px; }
        input { width: 100%; padding: 0.5rem; margin-bottom: 1rem; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
        button { width: 100%; padding: 0.5rem; background: #0070f3; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: #0051a2; }
        h2 { margin-top: 0; text-align: center; color: #333; }
    </style>
</head>
<body>
    <div class="login-box">
        <h2>Login</h2>
        <form action="/login" method="POST">
            <input type="hidden" name="return_to" value="` + returnTo + `">
            <div>
                <label>Username</label>
                <input type="text" name="username" required>
            </div>
            <div>
                <label>Password</label>
                <input type="password" name="password" required>
            </div>
            <button type="submit">Sign In</button>
        </form>
    </div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
