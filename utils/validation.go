package utils

import (
	"net/url"
	"strings"
)

// IsValidRedirect checks if the redirect URL is safe to redirect to.
// Strictly allows only relative paths to prevent Open Redirect vulnerabilities.
// It must start with / but not // (which is protocol relative).
func IsValidRedirect(redirectURL string) bool {
	u, err := url.Parse(redirectURL)
	if err != nil {
		return false
	}

	// Must be relative path
	if u.IsAbs() {
		return false
	}

	// Must start with /
	if !strings.HasPrefix(redirectURL, "/") {
		return false
	}

	// Must not be protocol relative
	if strings.HasPrefix(redirectURL, "//") {
		return false
	}

	// Must not contain control characters (CR/LF) - url.Parse usually handles this but good to be safe if we were doing manual parsing
	// but here we trust url.Parse failed if it had issues, or the string check.

	return true
}
