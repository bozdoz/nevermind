package common

import (
	"os"
	"path/filepath"
)

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
	return getNodeBinWithGetter(os.UserHomeDir, statFile, version, bin)
}

// test extraction (am I crazy?)
func getNodeBinWithGetter(homeGetter homeFunc, checkExists existsFunc, version Version, bin string) (path string, err error) {
	path, err = getNVMDirWithGetter(homeGetter, "node", string(version), "bin", bin)

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
