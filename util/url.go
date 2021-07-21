package util

import "net/url"

func Domain(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	return u.Scheme + "://" + u.Host
}