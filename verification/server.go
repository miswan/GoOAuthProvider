package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"html/template"
)

type HTMLHandler struct{}

func (h *HTMLHandler) Login(c echo.Context) error {
	next := c.QueryParam("next")

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Login - OAuth2 Provider</title>
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

func main() {
	e := echo.New()
	h := &HTMLHandler{}
	e.GET("/login", h.Login)
	e.Start(":8081")
}
