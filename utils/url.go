package utils

import (
	"net/url"
)

func BuildRedirectURL(baseURI string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURI)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}
