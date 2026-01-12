package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Index(c echo.Context) error {
	htmlContent := `
<!DOCTYPE html>
<html>
<head>
    <title>OAuth2 Provider</title>
    <style>
        body { font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; }
        .status { padding: 10px; background: #eef; border-radius: 5px; }
    </style>
</head>
<body>
    <h1>OAuth2 Provider</h1>
    <div class="status">
        <p>System is operational.</p>
        <p><a href="/login">Login</a></p>
    </div>
</body>
</html>`
	return c.HTML(http.StatusOK, htmlContent)
}

func (h *HTMLHandler) Login(c echo.Context) error {
	next := c.QueryParam("next")
	if next == "" {
		next = "/"
	}

	// Sanitize the next parameter to prevent XSS
	safeNext := html.EscapeString(next)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
    <style>
        body { font-family: sans-serif; max-width: 400px; margin: 0 auto; padding: 40px; }
        input { display: block; width: 100%%; margin-bottom: 10px; padding: 8px; }
        button { padding: 10px 20px; cursor: pointer; }
    </style>
</head>
<body>
    <h1>Login</h1>
    <form action="/login" method="POST">
        <input type="hidden" name="next" value="%s">
        <label>Username</label>
        <input type="text" name="username" required>
        <label>Password</label>
        <input type="password" name="password" required>
        <button type="submit">Login</button>
    </form>
    <p><a href="/register">Register</a></p>
</body>
</html>`, safeNext)
	return c.HTML(http.StatusOK, htmlContent)
}
