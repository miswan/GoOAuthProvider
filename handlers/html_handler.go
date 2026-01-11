package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HTMLHandler struct{}

func NewHTMLHandler() *HTMLHandler {
	return &HTMLHandler{}
}

func (h *HTMLHandler) Home(c echo.Context) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="OAuth2 Provider System Status and Documentation">
    <title>OAuth2 Provider - System Status</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 2rem;
            background-color: #fafafa;
        }
        main {
            background: white;
            padding: 2rem;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #2c3e50;
            border-bottom: 2px solid #eee;
            padding-bottom: 0.5rem;
            margin-top: 0;
        }
        .status {
            background: #e8f5e9;
            color: #2e7d32;
            padding: 1rem;
            border-radius: 6px;
            display: inline-flex;
            align-items: center;
            font-weight: bold;
            margin-bottom: 2rem;
            border: 1px solid #c8e6c9;
        }
        .endpoints {
            background: #f8f9fa;
            padding: 1.5rem;
            border-radius: 6px;
            border: 1px solid #e9ecef;
        }
        code {
            background: #eef1f6;
            padding: 0.2rem 0.4rem;
            border-radius: 4px;
            font-family: SFMono-Regular, Consolas, "Liberation Mono", Menlo, monospace;
            color: #24292e;
        }
        ul {
            list-style-type: none;
            padding: 0;
            margin: 0;
        }
        li {
            margin-bottom: 0.8rem;
            padding-left: 0.5rem;
            border-left: 3px solid #0366d6;
        }
        a {
            color: #0366d6;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
        footer {
            margin-top: 2rem;
            text-align: center;
            color: #666;
            font-size: 0.9rem;
        }
    </style>
</head>
<body>
    <main>
        <h1>OAuth2 Provider</h1>

        <div class="status" role="status" aria-live="polite">
            <span aria-hidden="true" style="margin-right: 8px;">✅</span> System Operational
        </div>

        <section class="endpoints" aria-labelledby="endpoints-title">
            <h2 id="endpoints-title">Available Endpoints</h2>
            <ul>
                <li><code>POST /register</code> - Register new user</li>
                <li><code>POST /login</code> - User login</li>
                <li><code>GET /authorize</code> - OAuth2 Authorization</li>
                <li><code>POST /token</code> - OAuth2 Token Exchange</li>
                <li><code>GET /client/:id</code> - Get Client Info</li>
            </ul>
        </section>
    </main>
    <footer>
        <p>Managed by Palette 🎨 &bull; <a href="https://github.com/oauth2-provider" aria-label="Project Repository">Documentation</a></p>
    </footer>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
