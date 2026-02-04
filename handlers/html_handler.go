package handlers

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
	"oauth2-provider/utils"
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
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #f3f4f6; }
        .card { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1); text-align: center; }
        h1 { color: #1f2937; margin-top: 0; }
        p { color: #4b5563; }
        .status { display: inline-flex; align-items: center; color: #059669; font-weight: 500; }
        .dot { width: 8px; height: 8px; background-color: #10b981; border-radius: 50%; margin-right: 8px; }
    </style>
</head>
<body>
    <div class="card">
        <h1>OAuth2 Provider</h1>
        <div class="status"><div class="dot"></div>System Operational</div>
        <p>This is a demonstration OAuth2 provider service.</p>
    </div>
</body>
</html>
`)
}

func (h *HTMLHandler) Login(c echo.Context) error {
	returnTo := c.QueryParam("return_to")
	if returnTo != "" && !utils.IsValidReturnTo(returnTo) {
		returnTo = ""
	}

	tmpl := template.Must(template.New("login").Parse(loginHTML))
	return tmpl.Execute(c.Response().Writer, map[string]interface{}{
		"ReturnTo": returnTo,
		"Error":    c.QueryParam("error"),
	})
}

const loginHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Login - OAuth2 Provider</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #f3f4f6; }
        .card { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1); width: 100%; max-width: 400px; }
        h1 { margin-top: 0; text-align: center; color: #1f2937; }
        .form-group { margin-bottom: 1rem; }
        label { display: block; margin-bottom: 0.5rem; color: #374151; font-weight: 500; }
        input { width: 100%; padding: 0.75rem; border: 1px solid #d1d5db; border-radius: 6px; box-sizing: border-box; margin-top: 0.25rem; }
        input:focus { outline: none; border-color: #2563eb; ring: 2px solid #93c5fd; }
        button { width: 100%; padding: 0.75rem; background-color: #2563eb; color: white; border: none; border-radius: 6px; cursor: pointer; font-weight: 600; margin-top: 1rem; }
        button:hover { background-color: #1d4ed8; }
        .error { background-color: #fee2e2; color: #991b1b; padding: 0.75rem; border-radius: 6px; margin-bottom: 1rem; text-align: center; font-size: 0.875rem; }
        .footer { text-align: center; margin-top: 1.5rem; font-size: 0.875rem; color: #6b7280; }
    </style>
</head>
<body>
    <div class="card">
        <h1>Sign In</h1>
        {{if .Error}}
            <div class="error">{{.Error}}</div>
        {{end}}
        <form action="/login" method="POST">
            {{if .ReturnTo}}
                <input type="hidden" name="return_to" value="{{.ReturnTo}}">
            {{end}}
            <div class="form-group">
                <label for="username">Username</label>
                <input type="text" id="username" name="username" required autofocus autocomplete="username">
            </div>
            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" name="password" required autocomplete="current-password">
            </div>
            <button type="submit">Sign In</button>
        </form>
        <div class="footer">
            OAuth2 Provider Service
        </div>
    </div>
</body>
</html>
`
