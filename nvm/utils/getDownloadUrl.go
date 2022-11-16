package utils

import (
	"fmt"
	"os"

	common "github.com/bozdoz/nevermind/nvm-common"
)

// Enum for downloading either node or shasums
type NodeDownload int

const (
	DOWNLOAD_NODE   NodeDownload = iota // downloads node
	DOWNLOAD_SHASUM                     // downloads SHASUMS
)

// makes the enum exclusive
type exclusive_node_file interface {
	x()
}

// ignore me
func (d NodeDownload) x() {}

// remote filename which indicates sha's for each downloadable file
const SHASUMS = "SHASUMS256.txt"

// default download location
const DEFAULT_BASE_URL = "https://nodejs.org/dist"

// NVM_NODEJS_ORG_MIRROR - base url for downloading node; default: "https://nodejs.org/dist"
var BASE_URL = os.Getenv("NVM_NODEJS_ORG_MIRROR")

func init() {
	// sets default for BASE_URL
	if BASE_URL == "" {
		BASE_URL = DEFAULT_BASE_URL
	}
}

/*
determine where to download node install files

`v` a verified/parsed version string

`d` is a [commands.NodeDownload]

used by [commands.Install]
*/
func GetDownloadUrl(v common.Version, d exclusive_node_file) string {
	remote_os, remote_arch := GetOsAndArch()
	ext := "tar.gz"

	// TODO: we don't actually test extracting zips on windows
	if remote_os == "win" {
		ext = "zip"
		panic("we haven't written any code that would extract a zip yet")
	}

	switch d {
	case DOWNLOAD_NODE:
		return fmt.Sprintf("%s/v%s/node-v%s-%s-%s.%s", BASE_URL, v, v, remote_os, remote_arch, ext)
	case DOWNLOAD_SHASUM:
		return fmt.Sprintf("%s/v%s/%s", BASE_URL, v, SHASUMS)
	default:
		return ""
	}
}
