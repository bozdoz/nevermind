package common

import (
	"fmt"
	"regexp"
)

// case-insensitive, leading v, major.minor.patch
func GetVersion(v string) (string, error) {
	re := regexp.MustCompile(`(?i)^v?(\d+(?:\.?\d+){0,2})$`)
	version := re.FindStringSubmatch(v)

	if len(version) < 2 {
		return "", fmt.Errorf("could not determine version: %s", v)
	}

	return version[1], nil
}
