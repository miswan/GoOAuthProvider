package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Home serves the welcome page for the API
func Home(c echo.Context) error {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OAuth2 Provider API</title>
    <style>
        :root {
            --primary: #3b82f6;
            --bg: #f8fafc;
            --text: #1e293b;
            --card-bg: #ffffff;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: var(--bg);
            color: var(--text);
            line-height: 1.5;
            margin: 0;
            padding: 2rem;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
        }
        main {
            background: var(--card-bg);
            padding: 2rem;
            border-radius: 8px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
            max-width: 600px;
            width: 100%;
        }
        h1 {
            color: var(--primary);
            margin-top: 0;
        }
        .status {
            display: inline-block;
            padding: 0.25rem 0.75rem;
            background-color: #dcfce7;
            color: #166534;
            border-radius: 9999px;
            font-size: 0.875rem;
            font-weight: 500;
            margin-bottom: 1.5rem;
        }
        ul {
            list-style: none;
            padding: 0;
        }
        li {
            padding: 0.75rem 0;
            border-bottom: 1px solid #e2e8f0;
            display: flex;
            align-items: center;
        }
        li:last-child {
            border-bottom: none;
        }
        .method {
            font-family: monospace;
            background: #e2e8f0;
            padding: 0.2rem 0.5rem;
            border-radius: 4px;
            font-size: 0.875rem;
            margin-right: 1rem;
            min-width: 60px;
            text-align: center;
        }
        .method.post { background: #dbeafe; color: #1e40af; }
        .method.get { background: #dcfce7; color: #166534; }
    </style>
</head>
<body>
    <main>
        <h1>OAuth2 Provider API</h1>
        <div class="status" role="status">● Systems Operational</div>
        <p>Welcome to the OAuth2 Provider API. Use the endpoints below to interact with the service.</p>

        <h2>Available Endpoints</h2>
        <ul>
            <li><span class="method post">POST</span> <span>/register</span></li>
            <li><span class="method post">POST</span> <span>/login</span></li>
            <li><span class="method get">GET</span> <span>/authorize</span></li>
            <li><span class="method post">POST</span> <span>/token</span></li>
            <li><span class="method get">GET</span> <span>/userinfo</span></li>
            <li><span class="method post">POST</span> <span>/client/register</span></li>
        </ul>
    </main>
</body>
</html>
`
	return c.HTML(http.StatusOK, html)
}
