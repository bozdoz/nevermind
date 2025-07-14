package common

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// TODO: monkey-patch this instead
type homeFunc func() (string, error)
type existsFunc func(path string) error

// user-relative path to app directory
const NVM_DIR string = ".nevermind"

// get a path relative to [common.NVM_DIR]
// returns directory string, and error
func GetNVMDir(path ...string) (string, error) {
	return getNVMDirWithGetter(os.UserHomeDir, path...)
}

// extracted for testing
func getNVMDirWithGetter(homeGetter homeFunc, path ...string) (string, error) {
	homeDir, err := homeGetter()

	path = append([]string{
		homeDir,
		NVM_DIR,
	},
		path...,
	)

	return filepath.Join(path...), err
}

// bin should likely be "node", "npm", "npx"; or any other node bin
// installed globally (e.g. yarn, typescript)
// returns directory string, and error
func GetNodeBin(version Version, bin string) (string, error) {
	found, err := getNodeBinWithGetter(os.UserHomeDir, statFile, version, bin)

	log.Printf("Looking for: %s\n\n  exists? %t\n\n", found, err == nil)

	return found, err
}

// test extraction (am I crazy?)
func getNodeBinWithGetter(homeGetter homeFunc, checkExists existsFunc, version Version, bin string) (path string, err error) {
	os := runtime.GOOS

	if os == "windows" {
		// maybe .exe?
		path, err = getNVMDirWithGetter(homeGetter, "node", string(version), fmt.Sprintf("%s.exe", bin))
	} else {
		path, err = getNVMDirWithGetter(homeGetter, "node", string(version), "bin", bin)
	}

	if err != nil {
		return
	}

	return path, checkExists(path)
}

// test extraction
func statFile(path string) error {
	_, err := os.Stat(path)

	return err
}
