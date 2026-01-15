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
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OAuth2 Provider Status</title>
    <style>
        :root { --bg: #111; --text: #eee; --accent: #3b82f6; --success: #22c55e; --card: #222; }
        body { background: var(--bg); color: var(--text); font-family: system-ui, -apple-system, sans-serif; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; }
        .card { background: var(--card); padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.3); max-width: 400px; width: 100%; }
        h1 { margin-top: 0; font-size: 1.5rem; display: flex; align-items: center; gap: 0.5rem; }
        .status { width: 10px; height: 10px; background: var(--success); border-radius: 50%; display: inline-block; box-shadow: 0 0 8px var(--success); }
        .links { margin-top: 1.5rem; display: flex; flex-direction: column; gap: 0.5rem; }
        a { color: var(--accent); text-decoration: none; padding: 0.5rem; border-radius: 4px; background: rgba(59, 130, 246, 0.1); text-align: center; transition: background 0.2s; }
        a:hover { background: rgba(59, 130, 246, 0.2); }
    </style>
</head>
<body>
    <div class="card">
        <h1><span class="status" aria-label="System Online"></span> System Online</h1>
        <p>OAuth2 Provider Service is running.</p>
        <nav class="links">
            <a href="/login">Login</a>
            <a href="/authorize?client_id=test&response_type=code&redirect_uri=http://localhost:8080/callback">Authorize (Test)</a>
        </nav>
    </div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}

func (h *HTMLHandler) Login(c echo.Context) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login - OAuth2 Provider</title>
    <style>
        :root { --bg: #111; --text: #eee; --accent: #3b82f6; --card: #222; --input: #333; }
        body { background: var(--bg); color: var(--text); font-family: system-ui, -apple-system, sans-serif; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; }
        .card { background: var(--card); padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.3); max-width: 350px; width: 100%; }
        h1 { margin-top: 0; font-size: 1.5rem; text-align: center; }
        form { display: flex; flex-direction: column; gap: 1rem; }
        label { font-size: 0.875rem; font-weight: 500; margin-bottom: 0.25rem; display: block; }
        input { background: var(--input); border: 1px solid #444; color: white; padding: 0.75rem; border-radius: 4px; font-size: 1rem; width: 100%; box-sizing: border-box; }
        input:focus { border-color: var(--accent); outline: none; ring: 2px solid var(--accent); }
        button { background: var(--accent); color: white; border: none; padding: 0.75rem; border-radius: 4px; font-size: 1rem; font-weight: 600; cursor: pointer; transition: opacity 0.2s; }
        button:hover { opacity: 0.9; }
        #message { margin-top: 1rem; text-align: center; font-size: 0.875rem; min-height: 1.25em; }
        .error { color: #ef4444; }
        .success { color: #22c55e; }
        a { color: #9ca3af; font-size: 0.875rem; text-decoration: none; text-align: center; display: block; margin-top: 1rem; }
    </style>
</head>
<body>
    <div class="card">
        <h1>Login</h1>
        <form id="loginForm">
            <div>
                <label for="username">Username</label>
                <input type="text" id="username" name="username" required autocomplete="username">
            </div>
            <div>
                <label for="password">Password</label>
                <input type="password" id="password" name="password" required autocomplete="current-password">
            </div>
            <button type="submit">Sign In</button>
        </form>
        <div id="message" role="alert"></div>
        <a href="/">← Back to Status</a>
    </div>
    <script>
        document.getElementById('loginForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const form = e.target;
            const msg = document.getElementById('message');
            msg.textContent = 'Authenticating...';
            msg.className = '';

            const formData = new URLSearchParams(new FormData(form));

            try {
                const res = await fetch('/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                    body: formData
                });

                const data = await res.json();

                if (res.ok) {
                    msg.textContent = 'Success! ' + (data.message || 'Logged in.');
                    msg.className = 'success';
                } else {
                    msg.textContent = data.message || 'Login failed';
                    msg.className = 'error';
                }
            } catch (err) {
                msg.textContent = 'Network error';
                msg.className = 'error';
            }
        });
    </script>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
