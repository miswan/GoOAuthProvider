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
    <title>OAuth2 Provider</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif; max-width: 800px; margin: 40px auto; padding: 20px; line-height: 1.6; color: #333; }
        h1 { color: #111; display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
        .status { display: inline-flex; align-items: center; padding: 4px 8px; border-radius: 4px; background: #e6fffa; color: #047481; font-weight: bold; font-size: 0.6em; vertical-align: middle; }
        .endpoints { background: #f7f9fa; padding: 20px; border-radius: 8px; margin-top: 20px; border: 1px solid #e1e4e8; }
        code { background: #eee; padding: 2px 4px; border-radius: 3px; font-family: monospace; font-size: 0.9em; }
        ul { list-style-type: none; padding: 0; }
        li { margin-bottom: 12px; padding-bottom: 8px; border-bottom: 1px solid #eee; }
        li:last-child { border-bottom: none; }
        .method { font-weight: bold; margin-right: 8px; min-width: 60px; display: inline-block; }
        .method.GET { color: #005cc5; }
        .method.POST { color: #22863a; }
        a { color: #0366d6; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <header>
        <h1>
            OAuth2 Provider
            <span class="status" role="status" aria-label="System Status: Operational">
                <span aria-hidden="true">●</span> Operational
            </span>
        </h1>
    </header>
    <main>
        <p>Welcome to the OAuth2 Provider service. This service provides authentication and authorization for client applications.</p>

        <section class="endpoints" aria-labelledby="endpoints-title">
            <h2 id="endpoints-title">Available Endpoints</h2>
            <ul>
                <li><span class="method GET">GET</span> <code>/authorize</code> - Authorization endpoint</li>
                <li><span class="method POST">POST</span> <code>/token</code> - Token exchange endpoint</li>
                <li><span class="method GET">GET</span> <code>/userinfo</code> - User information (Requires Auth)</li>
                <li><span class="method POST">POST</span> <code>/register</code> - User registration</li>
                <li><span class="method POST">POST</span> <code>/login</code> - User login</li>
                <li><span class="method POST">POST</span> <code>/client/register</code> - Client registration</li>
            </ul>
        </section>
    </main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
