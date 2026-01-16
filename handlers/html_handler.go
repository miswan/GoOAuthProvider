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
	return c.HTML(http.StatusOK, "<h1>System Operational</h1>")
}

func (h *HTMLHandler) Login(c echo.Context) error {
	next := c.QueryParam("next")
	html := `
<!DOCTYPE html>
<html>
<head><title>Login</title></head>
<body>
<h2>Login</h2>
<form action="/login" method="post">
  <input type="hidden" name="next" value="` + next + `">
  <label for="username">Username:</label><br>
  <input type="text" id="username" name="username"><br>
  <label for="password">Password:</label><br>
  <input type="password" id="password" name="password"><br><br>
  <input type="submit" value="Login">
</form>
</body>
</html>
`
	return c.HTML(http.StatusOK, html)
}
