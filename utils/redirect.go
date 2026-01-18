package utils

import (
	"net/url"
	"strings"
)

// IsValidRedirect checks if the redirect URL is valid.
// It allows only relative paths (starting with / but not //) to prevent open redirects.
func IsValidRedirect(redirectURL string) bool {
	if redirectURL == "" {
		return false
	}

	u, err := url.Parse(redirectURL)
	if err != nil {
		return false
	}

	// Must be relative path
	if u.Scheme != "" || u.Host != "" {
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

	return true
}
