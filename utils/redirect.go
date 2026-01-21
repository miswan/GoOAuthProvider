package utils

import (
	"strings"
)

func IsValidRedirect(url string) bool {
	return strings.HasPrefix(url, "/") && !strings.HasPrefix(url, "//")
}
