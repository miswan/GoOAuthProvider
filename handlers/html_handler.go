package handlers

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
	"strings"
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
			<style>body { font-family: sans-serif; text-align: center; padding: 50px; }</style>
		</head>
		<body>
			<h1>System Operational</h1>
			<p>The OAuth2 Provider is running.</p>
			<p><a href="/login">Login</a></p>
		</body>
		</html>
	`)
}

func (h *HTMLHandler) Login(c echo.Context) error {
	returnTo := c.QueryParam("return_to")

	// Open Redirect Protection
	if returnTo != "" && (!strings.HasPrefix(returnTo, "/") || strings.HasPrefix(returnTo, "//")) {
		returnTo = "/"
	}

	const loginTmpl = `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Login</title>
			<style>
				body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background: #f0f2f5; }
				.card { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); width: 300px; }
				input { width: 100%; padding: 10px; margin: 10px 0; box-sizing: border-box; border: 1px solid #ccc; border-radius: 4px; }
				button { width: 100%; padding: 10px; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
				button:hover { background: #0056b3; }
                .error { color: red; font-size: 0.9em; display: none; }
			</style>
		</head>
		<body>
			<div class="card">
				<h2>Login</h2>
                <div id="error" class="error"></div>
				<form id="loginForm">
					<input type="text" name="username" placeholder="Username" required>
					<input type="password" name="password" placeholder="Password" required>
                    <input type="hidden" name="return_to" value="{{.ReturnTo}}">
					<button type="submit">Sign In</button>
				</form>
			</div>
            <script>
                document.getElementById('loginForm').addEventListener('submit', async (e) => {
                    e.preventDefault();
                    const formData = new FormData(e.target);
                    const data = Object.fromEntries(formData.entries());
                    const returnTo = data.return_to;
                    delete data.return_to;

                    try {
                        const response = await fetch('/login', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify(data)
                        });

                        if (response.ok) {
                            if (returnTo) {
                                window.location.href = returnTo;
                            } else {
                                window.location.href = '/';
                            }
                        } else {
                            const result = await response.json();
                            const errorDiv = document.getElementById('error');
                            errorDiv.textContent = result.message || 'Login failed';
                            errorDiv.style.display = 'block';
                        }
                    } catch (err) {
                        console.error(err);
                    }
                });
            </script>
		</body>
		</html>
	`

	t, err := template.New("login").Parse(loginTmpl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	data := struct {
		ReturnTo string
	}{
		ReturnTo: returnTo,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.HTML(http.StatusOK, buf.String())
}
