package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// basename of binary in NVM_DIR
const NVM_SHIM = "nvm-shim"

// create symlinks of bins
func CreateSymlinks(version Version) (err error) {
	nvmDir, err := GetNVMDir()

	if err != nil {
		return
	}

	// directory where we create the symlinks
	symlinkDir := filepath.Join(nvmDir, "bin")

	// target all symlinks to nvmShim
	nvmShim := filepath.Join(symlinkDir, NVM_SHIM)
	// installed bins (via npm i, or node install)
	nodeBins := filepath.Join(nvmDir, "node", string(version), "bin")

	bins, err := os.ReadDir(nodeBins)

	if err != nil {
		return
	}

	for _, bin := range bins {
		symlink := filepath.Join(symlinkDir, bin.Name())

		err = syscall.Symlink(nvmShim, symlink)

		if err != nil && !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("cannot make symlink for %s = %w", bin.Name(), err)
		}
	}

	return
}
