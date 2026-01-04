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

func (h *HTMLHandler) Login(c echo.Context) error {
	redirectTo := c.QueryParam("redirect_to")
	if redirectTo == "" {
		redirectTo = "/"
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Login - OAuth2 Provider</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; }
        .login-container { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); width: 300px; }
        h2 { margin-top: 0; color: #333; }
        input { width: 100%%; padding: 10px; margin: 10px 0; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
        button { width: 100%%; padding: 10px; background-color: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background-color: #0056b3; }
        .error { color: red; font-size: 0.9em; display: none; }
    </style>
</head>
<body>
    <div class="login-container">
        <h2>Login</h2>
        <div id="error-msg" class="error"></div>
        <form id="login-form">
            <input type="text" name="username" placeholder="Username" required>
            <input type="password" name="password" placeholder="Password" required>
            <button type="submit">Sign In</button>
        </form>
        <p style="font-size: 0.8em; text-align: center; margin-top: 10px;">
            Don't have an account? <a href="/register">Register</a>
        </p>
    </div>
    <script>
        document.getElementById('login-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const data = Object.fromEntries(formData.entries());

            try {
                const response = await fetch('/login?redirect_to=%s', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(data)
                });

                if (response.ok) {
                    const result = await response.json();
                    window.location.href = result.redirect_to || '%s';
                } else {
                    const err = await response.json();
                    const errorDiv = document.getElementById('error-msg');
                    errorDiv.textContent = err.message || 'Login failed';
                    errorDiv.style.display = 'block';
                }
            } catch (error) {
                console.error('Error:', error);
            }
        });
    </script>
</body>
</html>`, redirectTo, redirectTo)

	return c.HTML(http.StatusOK, html)
}

func (h *HTMLHandler) Register(c echo.Context) error {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Register - OAuth2 Provider</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; }
        .login-container { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); width: 300px; }
        h2 { margin-top: 0; color: #333; }
        input { width: 100%; padding: 10px; margin: 10px 0; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
        button { width: 100%; padding: 10px; background-color: #28a745; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background-color: #218838; }
        .error { color: red; font-size: 0.9em; display: none; }
        .success { color: green; font-size: 0.9em; display: none; }
    </style>
</head>
<body>
    <div class="login-container">
        <h2>Register</h2>
        <div id="error-msg" class="error"></div>
        <div id="success-msg" class="success"></div>
        <form id="register-form">
            <input type="text" name="username" placeholder="Username" required>
            <input type="email" name="email" placeholder="Email" required>
            <input type="password" name="password" placeholder="Password" required>
            <button type="submit">Register</button>
        </form>
        <p style="font-size: 0.8em; text-align: center; margin-top: 10px;">
            Already have an account? <a href="/login">Login</a>
        </p>
    </div>
    <script>
        document.getElementById('register-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const data = Object.fromEntries(formData.entries());

            try {
                const response = await fetch('/register', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(data)
                });

                if (response.ok) {
                    document.getElementById('success-msg').textContent = 'Registration successful! Redirecting to login...';
                    document.getElementById('success-msg').style.display = 'block';
                    document.getElementById('error-msg').style.display = 'none';
                    setTimeout(() => window.location.href = '/login', 2000);
                } else {
                    const err = await response.json();
                    document.getElementById('error-msg').textContent = err.message || 'Registration failed';
                    document.getElementById('error-msg').style.display = 'block';
                    document.getElementById('success-msg').style.display = 'none';
                }
            } catch (error) {
                console.error('Error:', error);
            }
        });
    </script>
</body>
</html>`

	return c.HTML(http.StatusOK, html)
}

func (h *HTMLHandler) Home(c echo.Context) error {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>OAuth2 Provider</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; margin: 0; }
        .container { text-align: center; background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #333; }
        .links { margin-top: 20px; }
        a { text-decoration: none; color: #007bff; margin: 0 10px; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <div class="container">
        <h1>OAuth2 Provider</h1>
        <p>System is running.</p>
        <div class="links">
            <a href="/login">Login</a>
            <a href="/register">Register</a>
        </div>
    </div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
