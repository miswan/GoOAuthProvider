package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Home(c echo.Context) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OAuth2 Provider</title>
    <style>
        :root {
            --primary-color: #2563eb;
            --text-color: #1f2937;
            --bg-color: #f3f4f6;
            --card-bg: #ffffff;
        }
        body {
            font-family: system-ui, -apple-system, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            line-height: 1.5;
            background-color: var(--bg-color);
            color: var(--text-color);
            margin: 0;
            padding: 2rem;
            display: flex;
            justify-content: center;
            min-height: 100vh;
        }
        main {
            background-color: var(--card-bg);
            padding: 2.5rem;
            border-radius: 0.75rem;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
            max-width: 600px;
            width: 100%;
        }
        h1 {
            color: #111827;
            margin-top: 0;
            font-size: 1.875rem;
            letter-spacing: -0.025em;
        }
        p {
            color: #4b5563;
            margin-bottom: 1.5rem;
        }
        .status {
            display: inline-flex;
            align-items: center;
            padding: 0.5rem 1rem;
            background-color: #dcfce7;
            color: #166534;
            border-radius: 9999px;
            font-weight: 500;
            font-size: 0.875rem;
        }
        .status svg {
            margin-right: 0.5rem;
            width: 1.25rem;
            height: 1.25rem;
        }
        .actions {
            margin-top: 2rem;
            border-top: 1px solid #e5e7eb;
            padding-top: 1.5rem;
        }
        .link {
            color: var(--primary-color);
            text-decoration: none;
            font-weight: 500;
        }
        .link:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <main>
        <h1>OAuth2 Provider</h1>
        <p>Welcome to the OAuth2 Provider Service. This secure service manages user authentication and client authorization flows.</p>

        <div class="status" role="status">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            System is operational
        </div>

        <div class="actions">
            <p>Ready to integrate?</p>
            <a href="/authorize" class="link">Start Authorization Flow</a>
        </div>
    </main>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
