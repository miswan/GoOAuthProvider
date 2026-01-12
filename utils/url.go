package utils

import (
	"net/url"
	"strings"
)

// IsValidRedirect checks if the redirect URL is safe (relative path).
func IsValidRedirect(redirectURL string) bool {
	return strings.HasPrefix(redirectURL, "/") && !strings.HasPrefix(redirectURL, "//")
}

// BuildRedirectURL constructs a URL with the given query parameters.
func BuildRedirectURL(baseURL string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}
