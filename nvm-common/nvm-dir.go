package common

import (
	"os"
	"path/filepath"
)

// user-relative path to app directory
const NVM_DIR string = ".nevermind"

// get a path relative to [common.NVM_DIR]
func GetNVMDir(path ...string) string {
	// TODO: deal with this error
	homeDir, _ := os.UserHomeDir()

	path = append([]string{
		homeDir,
		NVM_DIR,
	},
		path...,
	)

	return filepath.Join(path...)
}

// bin should likely be "node", "npm", "npx"; or any other node bin
// installed globally (e.g. yarn, typescript)
func GetNodeBin(version Version, bin string) (dir string, err error) {
	dir = GetNVMDir("node", string(version), "bin", bin)

	_, err = os.Stat(dir)

	return
}
