package utils

import (
	"errors"
	"io"
	"os"

	common "github.com/bozdoz/nevermind/nvm-common"
)

var ErrNoNvmRc = errors.New("no .nvmrc found")

// get version in .nvmrc
func ReadNvmrc() (ver common.Version, err error) {
	// check for .nvmrc
	_, err = os.Stat(".nvmrc")

	if err != nil {
		err = ErrNoNvmRc
		return
	}

	// check for version in .nvmrc
	file, err := os.Open(".nvmrc")

	if err != nil {
		return
	}

	defer file.Close()

	body, err := io.ReadAll(file)

	if err != nil {
		return
	}

	ver, err = common.GetVersion(string(body))

	return
}
