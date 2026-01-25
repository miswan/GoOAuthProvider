package utils

import "strings"

func IsValidRedirect(url string) bool {
	if len(url) == 0 {
		return false
	}
	// Must start with /
	if !strings.HasPrefix(url, "/") {
		return false
	}
	// Must NOT start with // (protocol relative)
	if strings.HasPrefix(url, "//") {
		return false
	}
	return true
}
