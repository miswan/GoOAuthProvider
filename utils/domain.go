package utils

import (
	"net/url"
	"strings"
)

// IsValidRedirect checks if the urlStr is safe to redirect to.
// It allows relative paths (starting with / but not //)
// It allows absolute URLs if they match the host.
func IsValidRedirect(host, urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check for dangerous schemes
	if u.Scheme != "" && u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	// Check if it's a relative path
	if u.Host == "" {
		// Prevent protocol-relative URLs (e.g. //evil.com) which Parse treats as relative path with empty host?
		// Actually url.Parse("//evil.com") results in Host="evil.com" but Scheme=""
		// So checking Host == "" is not enough if we don't check path content.

		// Let's check the string itself to be sure
		if strings.HasPrefix(urlStr, "//") {
			return false
		}

		// Must start with / if it's relative
		if !strings.HasPrefix(urlStr, "/") {
			return false
		}

		return true
	}

	// If it has a host, it must match our host
	// Normalize hosts (remove port if any)
	h := host
	if strings.Contains(h, ":") {
		h = strings.Split(h, ":")[0]
	}

	uHost := u.Host
	if strings.Contains(uHost, ":") {
		uHost = strings.Split(uHost, ":")[0]
	}

	return strings.EqualFold(h, uHost)
}

// Keeping IsSameDomain for backward compatibility if needed, but implementation updated
func IsSameDomain(host, urlStr string) bool {
	return IsValidRedirect(host, urlStr)
}
