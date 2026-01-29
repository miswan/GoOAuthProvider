package handlers

import (
	"encoding/json"
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
	returnTo := c.QueryParam("return_to")
	returnToJSON, _ := json.Marshal(returnTo)
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; background: #111; color: #fff; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
        .container { background: #222; padding: 2rem; border-radius: 8px; width: 300px; }
        input { display: block; width: 100%; margin-bottom: 1rem; padding: 0.5rem; background: #333; border: 1px solid #444; color: #fff; border-radius: 4px; box-sizing: border-box; }
        button { width: 100%; padding: 0.5rem; background: #4ade80; color: #000; border: none; border-radius: 4px; cursor: pointer; font-weight: bold; }
        @media (prefers-reduced-motion: reduce) { * { animation: none !important; transition: none !important; } }
        @media (prefers-color-scheme: light) { body { background: #f0f0f0; color: #000; } .container { background: #fff; border: 1px solid #ccc; } input { background: #fff; border: 1px solid #ccc; color: #000; } }
    </style>
</head>
<body>
    <div class="container">
        <h2>Login</h2>
        <form id="loginForm">
            <input type="text" name="username" placeholder="Username" required>
            <input type="password" name="password" placeholder="Password" required>
            <button type="submit">Login</button>
        </form>
    </div>
    <script>
        document.getElementById('loginForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const data = Object.fromEntries(formData.entries());

            try {
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(data)
                });

                if (response.ok) {
                    const result = await response.json();
                    const returnTo = ` + string(returnToJSON) + `;
                    if (returnTo) {
                        window.location.href = returnTo;
                    } else {
                        alert('Login successful');
                    }
                } else {
                    alert('Login failed');
                }
            } catch (error) {
                console.error('Error:', error);
                alert('An error occurred');
            }
        });
    </script>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
