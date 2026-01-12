package utils

import (
	"testing"
)

func TestIsValidRedirect(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"/valid/path", true},
		{"/", true},
		{"//invalid/path", false},
		{"https://google.com", false},
		{"javascript:alert(1)", false},
		{"/path?query=1", true},
	}

	for _, tt := range tests {
		if got := IsValidRedirect(tt.url); got != tt.expected {
			t.Errorf("IsValidRedirect(%q) = %v; want %v", tt.url, got, tt.expected)
		}
	}
}

func TestBuildRedirectURL(t *testing.T) {
	baseURL := "https://client.com/callback"
	params := map[string]string{
		"code":  "123",
		"state": "xyz",
	}

	got, err := BuildRedirectURL(baseURL, params)
	if err != nil {
		t.Fatalf("BuildRedirectURL failed: %v", err)
	}

	// Order of query params is not guaranteed, so check contains
	if got != "https://client.com/callback?code=123&state=xyz" && got != "https://client.com/callback?state=xyz&code=123" {
		t.Errorf("BuildRedirectURL() = %q; want params", got)
	}
}
