package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

const commonStyles = `
<style>
:root { --primary: #2563eb; --bg: #f8fafc; --text: #1e293b; --surface: #ffffff; --border: #cbd5e1; }
@media (prefers-color-scheme: dark) {
  :root { --primary: #3b82f6; --bg: #0f172a; --text: #f1f5f9; --surface: #1e293b; --border: #334155; }
}
body { font-family: system-ui, -apple-system, sans-serif; background: var(--bg); color: var(--text); display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; line-height: 1.5; }
.card { background: var(--surface); padding: 2rem; border-radius: 1rem; box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1); width: 100%; max-width: 24rem; border: 1px solid var(--border); }
h1 { margin-top: 0; color: var(--primary); font-size: 1.5rem; margin-bottom: 1.5rem; }
input { display: block; width: 100%; padding: 0.75rem; margin-bottom: 1rem; border: 1px solid var(--border); border-radius: 0.375rem; background: var(--bg); color: var(--text); box-sizing: border-box; }
input:focus { outline: 2px solid var(--primary); outline-offset: 2px; border-color: transparent; }
button { background: var(--primary); color: white; border: none; padding: 0.75rem 1rem; border-radius: 0.375rem; cursor: pointer; width: 100%; font-weight: 600; font-size: 1rem; transition: opacity 0.2s; }
button:hover { opacity: 0.9; }
button:focus-visible { outline: 2px solid var(--primary); outline-offset: 2px; }
label { display: block; margin-bottom: 0.5rem; font-weight: 500; font-size: 0.875rem; }
.status { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.5rem; padding: 1rem; background: var(--bg); border-radius: 0.5rem; }
.dot { width: 0.75rem; height: 0.75rem; background: #22c55e; border-radius: 50%; box-shadow: 0 0 0 2px var(--bg); }
.link { color: var(--primary); text-decoration: none; font-size: 0.875rem; display: block; text-align: center; margin-top: 1rem; }
.link:hover { text-decoration: underline; }
</style>
`

func (h *HTMLHandler) Index(c echo.Context) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>System Status - OAuth2 Provider</title>
    ` + commonStyles + `
</head>
<body>
    <main class="card">
        <h1>OAuth2 Provider</h1>
        <div class="status" role="status">
            <div class="dot" aria-hidden="true"></div>
            <span>System Operational</span>
        </div>
        <p>This is a demonstration OAuth2 provider.</p>
        <a href="/login" class="button" style="display:block; text-align:center; background:var(--primary); color:white; text-decoration:none; padding:0.75rem; border-radius:0.375rem; font-weight:600;">Log In</a>
    </main>
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
    ` + commonStyles + `
</head>
<body>
    <main class="card">
        <h1>Welcome Back</h1>
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
        <div id="message" aria-live="polite" style="margin-top:1rem; text-align:center; display:none;"></div>
        <a href="/" class="link">Back to Home</a>
    </main>
    <script>
        document.getElementById('loginForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const btn = e.target.querySelector('button');
            const msg = document.getElementById('message');

            btn.disabled = true;
            btn.textContent = 'Signing in...';
            msg.style.display = 'none';

            try {
                const formData = new FormData(e.target);
                const data = Object.fromEntries(formData.entries());

                const res = await fetch('/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(data)
                });

                if (res.ok) {
                    msg.style.color = '#22c55e';
                    msg.textContent = 'Success! Redirecting...';
                    msg.style.display = 'block';
                    setTimeout(() => window.location.href = '/', 1000);
                } else {
                    const err = await res.json();
                    throw new Error(err.message || 'Login failed');
                }
            } catch (err) {
                msg.style.color = '#ef4444';
                msg.textContent = err.message;
                msg.style.display = 'block';
                btn.disabled = false;
                btn.textContent = 'Sign In';
            }
        });
    </script>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
