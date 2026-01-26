package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Login(c echo.Context) error {
	continueTo := c.QueryParam("continue_to")

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Login - OAuth2 Provider</title>
    <style>
        body { font-family: system-ui, -apple-system, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background: #f4f4f5; }
        .card { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1); width: 100%; max-width: 400px; }
        h1 { margin-top: 0; text-align: center; }
        input { display: block; width: 100%; padding: 0.5rem; margin-bottom: 1rem; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
        button { width: 100%; padding: 0.5rem; background: #2563eb; color: white; border: none; border-radius: 4px; cursor: pointer; font-weight: bold; }
        button:hover { background: #1d4ed8; }
        .error { color: #ef4444; margin-bottom: 1rem; text-align: center; display: none; }
    </style>
</head>
<body>
    <div class="card">
        <h1>Sign In</h1>
        <div id="error" class="error"></div>
        <form id="loginForm" action="/login" method="POST">
            <input type="hidden" name="continue_to" value="` + continueTo + `">
            <input type="text" name="username" placeholder="Username" required>
            <input type="password" name="password" placeholder="Password" required>
            <button type="submit">Login</button>
        </form>
		<p style="text-align: center; margin-top: 1rem;"><a href="/register">Create an account</a></p>
    </div>
	<script>
		// Handle errors from query param
		const urlParams = new URLSearchParams(window.location.search);
		const error = urlParams.get('error');
		if (error) {
			const errorDiv = document.getElementById('error');
			errorDiv.textContent = error;
			errorDiv.style.display = 'block';
		}
	</script>
</body>
</html>
`
	return c.HTML(http.StatusOK, html)
}
