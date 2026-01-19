package utils

import (
	"net/url"
	"strings"
)

// IsValidRedirect checks if the redirect URL is safe to use.
// It allows only relative paths (starting with / but not //) to prevent open redirects.
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

	// Parse to check for other issues
	u, err := url.Parse(redirectURL)
	if err != nil {
		return false
	}

	// Ensure no scheme or host is present (relative path only)
	if u.Scheme != "" || u.Host != "" {
		return false
	}

	return true
}
