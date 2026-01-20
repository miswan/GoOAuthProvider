package utils

import "strings"

func IsValidRedirect(url string) bool {
	// Only allow relative URLs to prevent open redirects
	// Must start with /
	// Must NOT start with // (protocol relative)
	return strings.HasPrefix(url, "/") && !strings.HasPrefix(url, "//")
}
