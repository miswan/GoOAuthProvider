package utils

import (
	"net/url"
	"strings"
)

func IsSameDomain(host, redirectURL string) bool {
	u, err := url.Parse(redirectURL)
	if err != nil {
		return false
	}
	// Strip port from host if present
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	// Strip port from u.Hostname() if present (it shouldn't be, but safe to check)
	redirectHost := u.Hostname()

	return strings.EqualFold(host, redirectHost)
}
