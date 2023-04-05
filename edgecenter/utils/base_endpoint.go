package utils

import (
	"net/url"
	"regexp"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func BaseRootEndpoint(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

// NormalizeURLPath removes duplicated slashes.
func NormalizeURLPath(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	path := u.Path
	r := regexp.MustCompile(`//+`)
	u.Path = r.ReplaceAllLiteralString(path, "/")
	return edgecloud.NormalizeURL(u.String()), nil
}
