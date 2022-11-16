package utils

import "runtime"

// TODO: this can be cached, or run once
// used to determine download url in [commands.GetDownloadUrl]
func GetOsAndArch() (remote_os, remote_arch string) {
	remote_arch = runtime.GOARCH

	switch remote_arch {
	case "x86_64", "amd64":
		remote_arch = "x64"
	case "aarch64":
		remote_arch = "arm64"
	}

	remote_os = runtime.GOOS

	if remote_os == "windows" {
		remote_os = "win"
		if remote_arch == "amd64" {
			remote_arch = "x64"
		} else {
			remote_arch = "x86"
		}
	}

	return
}
