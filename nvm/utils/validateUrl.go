package utils

import (
	"fmt"
	"net/url"
)

// check that a url string is valid
func ValidateUrl(url_str string) (string, error) {
	u, err := url.Parse(url_str)

	if err != nil {
		return "", fmt.Errorf("%s is not a valid url", url_str)
	}

	return u.String(), nil
}
