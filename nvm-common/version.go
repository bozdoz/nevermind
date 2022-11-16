package common

import (
	"fmt"
	"regexp"
)

// a parsed, verified version string
// returned by [common.GetVersion]
type Version string

// case-insensitive, leading v, major.minor.patch
func GetVersion(v string) (Version, error) {
	re := regexp.MustCompile(`(?i)^v?(\d+(?:\.?\d+){0,2})$`)
	version := re.FindStringSubmatch(v)

	if len(version) < 2 {
		return "", versionError(v)
	}

	return Version(version[1]), nil
}

func versionError(v string) error {
	return fmt.Errorf("could not determine version: %s", v)
}
