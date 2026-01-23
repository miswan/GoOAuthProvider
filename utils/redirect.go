package utils

import (
	"strings"
)

// IsValidRedirect checks if the redirect URI is a valid relative path.
// It allows paths starting with "/" but rejects "//" to prevent open redirects.
func IsValidRedirect(uri string) bool {
	return strings.HasPrefix(uri, "/") && !strings.HasPrefix(uri, "//")
}
