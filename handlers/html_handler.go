package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) LoginView(c echo.Context) error {
	redirectTo := c.QueryParam("redirect_to")
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
</head>
<body>
    <h2>Login</h2>
    <form action="/login" method="POST">
        <input type="hidden" name="redirect_to" value="%s">
        <label>Username: <input type="text" name="username" required></label><br>
        <label>Password: <input type="password" name="password" required></label><br>
        <button type="submit">Login</button>
    </form>
</body>
</html>
`, redirectTo)
	return c.HTML(http.StatusOK, html)
}
