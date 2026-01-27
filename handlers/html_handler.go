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
	return c.HTML(http.StatusOK, "<h1>System Operational</h1><a href='/login'>Login</a>")
}

func (h *HTMLHandler) Login(c echo.Context) error {
    continueTo := c.QueryParam("continue_to")
    if continueTo == "" {
        continueTo = "/"
    }

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f0f0; }
        .login-container { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); width: 300px; }
        h2 { text-align: center; color: #333; }
        input { width: 100%; padding: 0.5rem; margin-bottom: 1rem; border: 1px solid #ccc; border-radius: 4px; box-sizing: border-box; }
        button { width: 100%; padding: 0.5rem; background-color: #2563eb; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background-color: #1d4ed8; }
    </style>
</head>
<body>
    <div class="login-container">
        <h2>Login</h2>
        <form id="loginForm">
            <input type="text" id="username" name="username" placeholder="Username" required>
            <input type="password" id="password" name="password" placeholder="Password" required>
            <input type="hidden" id="continue_to" value="` + continueTo + `">
            <button type="submit">Sign In</button>
        </form>
        <div id="message" style="color: red; margin-top: 10px; text-align: center;"></div>
    </div>
    <script>
        document.getElementById('loginForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            const continueTo = document.getElementById('continue_to').value;

            try {
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username, password })
                });

                if (response.ok) {
                    window.location.href = continueTo;
                } else {
                    const data = await response.json();
                    document.getElementById('message').innerText = data.message || 'Login failed';
                }
            } catch (err) {
                document.getElementById('message').innerText = 'An error occurred';
            }
        });
    </script>
</body>
</html>
`
	return c.HTML(http.StatusOK, html)
}
