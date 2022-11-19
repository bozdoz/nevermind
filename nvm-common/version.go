package common

import (
	"fmt"
	"regexp"
	"strings"
)

// a parsed, verified version string
// returned by [common.GetVersion]
type Version string

// has all major.minor.patch
func (v Version) IsSpecific() bool {
	return strings.Count(string(v), ".") == 2
}

var version_regex = regexp.MustCompile(`(?i)^v?(\d+(?:\.?\d+){0,2})$`)

// case-insensitive, leading v, major.minor.patch
func GetVersion(v string) (Version, error) {
	version := version_regex.FindStringSubmatch(v)

	if len(version) < 2 {
		return "", versionError(v)
	}

	return Version(version[1]), nil
}

func versionError(v string) error {
	return fmt.Errorf("could not determine version: %s", v)
}
