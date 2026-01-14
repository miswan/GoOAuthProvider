package utils

import (
	"net/url"
	"strings"
)

// IsValidRedirect checks if the redirect URL is relative and safe
func IsValidRedirect(redirectURL string) bool {
	u, err := url.Parse(redirectURL)
	if err != nil {
		return false
	}
	// Only allow relative paths (start with /) and not // (protocol relative)
	return strings.HasPrefix(redirectURL, "/") && !strings.HasPrefix(redirectURL, "//") && u.Host == "" && u.Scheme == ""
}
