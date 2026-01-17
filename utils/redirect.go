package utils

import (
	"net/url"
	"strings"
)

// IsValidRedirect checks if the redirect URL is a valid relative path.
// It prevents open redirects by ensuring the path starts with '/' but not '//'.
func IsValidRedirect(redirectURL string) bool {
	if redirectURL == "" {
		return false
	}

	// Must start with /
	if !strings.HasPrefix(redirectURL, "/") {
		return false
	}

	// Must not start with // (protocol relative)
	if strings.HasPrefix(redirectURL, "//") {
		return false
	}

	// Parse to ensure it's valid
	u, err := url.Parse(redirectURL)
	if err != nil {
		return false
	}

	// Ensure no host is set (relative)
	if u.Scheme != "" || u.Host != "" {
		return false
	}

	return true
}
