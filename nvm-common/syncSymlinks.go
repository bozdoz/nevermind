package common

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

// basename of binary in NVM_DIR
const NVM_SHIM = "nvm-shim"
const NVM = "nvm"
const BIN = "bin"

var safeFiles = map[string]bool{
	NVM_SHIM: true,
	NVM:      true,
}

// sync symlinks of bins with our own bin directory
//
// needs to be wiped and re-established each time node version is changed
func SyncSymlinks(version Version) (err error) {
	nvmDir, err := GetNVMDir()

	if err != nil {
		log.Println("nvm dir err")
		return
	}

	// directory where we create the symlinks
	symlinkDir := filepath.Join(nvmDir, BIN)

	// target all symlinks to nvmShim
	nvmShim := filepath.Join(symlinkDir, NVM_SHIM)
	// installed bins (via npm i, or node install)
	nodeBins := filepath.Join(nvmDir, "node", version.String(), BIN)

	// create symlinks for all bins in node/{version}/bin
	bins, err := os.ReadDir(nodeBins)

	if err != nil {
		log.Println("node bin read err")
		return
	}

	added := map[string]bool{}

	for _, bin := range bins {
		name := bin.Name()
		symlink := filepath.Join(symlinkDir, name)

		err = syscall.Symlink(nvmShim, symlink)

		if err != nil && !errors.Is(err, os.ErrExist) {
			log.Println("symlink error")
			return fmt.Errorf("cannot make symlink for %s = %w", bin.Name(), err)
		}
		log.Println("symlinked", symlink)

		added[name] = true
	}

	return RemoveSymlinks(version, added)
}

// removes symlinks created by us to point to nvm-shim
//
// ignores any filenames found in `except`
func RemoveSymlinks(version Version, except map[string]bool) (err error) {
	symlinkDir, err := GetNVMDir(BIN)

	if err != nil {
		log.Println("symlink dir err")
		return
	}

	bins, err := os.ReadDir(symlinkDir)

	if err != nil {
		// could get into a weird state here (out-of-sync)
		// should we abort the whole process if there's any error, somehow?
		log.Println("read symlink dir err")
		return
	}

	for _, bin := range bins {
		name := bin.Name()

		if !except[name] && !safeFiles[name] {
			// not added and not safe means we delete
			symlink := filepath.Join(symlinkDir, name)
			log.Println("removing out of sync package:", symlink)
			err = os.Remove(symlink)

			if err != nil {
				log.Println("remove symlink err")
				return
			}
		}
	}

	return
}
