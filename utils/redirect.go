package utils

import (
	"net/url"
	"strings"
)

// IsValidReturnTo checks if the return_to URL is a valid relative path
// to prevent open redirect vulnerabilities.
func IsValidReturnTo(uri string) bool {
	if uri == "" {
		return false
	}

	u, err := url.Parse(uri)
	if err != nil {
		return false
	}

	// Must be relative (no scheme or host)
	if u.Scheme != "" || u.Host != "" {
		return false
	}

	// Prevent protocol-relative URLs (//example.com)
	if strings.HasPrefix(uri, "//") {
		return false
	}

	// Must start with /
	if !strings.HasPrefix(uri, "/") {
		return false
	}

	return true
}
