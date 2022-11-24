package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type VersionError struct {
	message string
	version string
}

func (ve VersionError) Error() string {
	return fmt.Sprintf("%s: %s", ve.message, ve.version)
}

// a semver string without the leading "v"
// returned by [common.GetVersion]
type Version string

func (v Version) String() string {
	return string(v)
}

// has all major.minor.patch
func (v Version) IsSpecific() bool {
	return strings.Count(string(v), ".") == 2
}

func (v Version) getIndex(i int) int {
	parts := strings.Split(string(v), ".")
	if len(parts) > i {
		str := parts[i]
		val, _ := strconv.Atoi(str)
		return val
	}
	return -1
}

func (v Version) Major() int {
	return v.getIndex(0)
}

func (v Version) Minor() int {
	return v.getIndex(1)
}

func (v Version) Patch() int {
	return v.getIndex(2)
}

var version_regex = regexp.MustCompile(`(?i)^v?(\d+(?:\.?\d+){0,2})$`)

// case-insensitive, leading v, major.minor.patch
func GetVersion(v string) (Version, error) {
	version := version_regex.FindStringSubmatch(v)

	if len(version) < 2 {
		// TODO: I'm not sure how to check if errors.Is(err, VersionError)
		return "", VersionError{"could not determine version", v}
	}

	return Version(version[1]), nil
}
