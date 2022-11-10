package common

import (
	"os"
	"path/filepath"
)

var NVMDIR = ".nevermind"

func GetNVMDir(path ...string) string {
	// TODO: deal with this error
	homeDir, _ := os.UserHomeDir()

	path = append([]string{
		homeDir,
		NVMDIR,
	},
		path...,
	)

	return filepath.Join(path...)
}

// bin should be "node", "npm", "npx"
func GetNodeBin(version, bin string) (dir string, err error) {
	dir = GetNVMDir("node", version, "bin", bin)

	_, err = os.Stat(dir)

	return
}
