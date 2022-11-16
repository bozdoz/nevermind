/*
commands holds all flagset commands for nvm:

  - install
  - use
  - list
  - uninstall
  - which

Example

	nvm install 18.0.0
	nvm use 18.0.0
	node -v // v18.0.0
*/
package commands

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/utils"
)

// basename of binary in NVM_DIR
const NVM_SHIM = "nvm-shim"

// TODO: document env variables that we use

// flag set for [commands.Install]
var InstallCmd = flag.NewFlagSet("install", flag.ContinueOnError)

var installHelp = InstallCmd.Bool("help", false, `prints help text`)
var installNoUse = InstallCmd.Bool("no-use", false, `disable calling nvm use <version> after install`)

func init() {
	InstallCmd.Usage = func() {
		utils.PrintTabs("\tinstall, i\tinstall a version of node")
	}
}

// install a given node version
// TODO: nvm install -h outputs help text twice
func Install(args []string) (err error) {
	InstallCmd.Parse(args)
	args = InstallCmd.Args()

	if *installHelp {
		// TODO: maybe verbose help message here
		InstallCmd.Usage()
		utils.FlushTabs()
		return
	}

	if len(args) == 0 {
		return errors.New("install did not get arguments")
	}
	version, err := common.GetVersion(args[0])

	if err != nil {
		return
	}

	log.Println("installing version", version)

	node_url := utils.GetDownloadUrl(version, utils.DOWNLOAD_NODE)
	sha_url := utils.GetDownloadUrl(version, utils.DOWNLOAD_SHASUM)
	log.Println("download node url", node_url)
	log.Println("download sha url", sha_url)

	node_chan := make(chan []byte)
	sha_chan := make(chan []byte)
	err_chan := make(chan error)

	// TODO: allow cancelling the downloads via SIGINT
	go utils.Download(node_url, node_chan, err_chan)
	go utils.Download(sha_url, sha_chan, err_chan)

	var node_body []byte
	var sha_body []byte

	select {
	case node_body = <-node_chan:
	case err = <-err_chan:
		return
	}

	select {
	case sha_body = <-sha_chan:
	case err = <-err_chan:
		return
	}

	segments := strings.Split(node_url, "/")
	baseFileName := segments[len(segments)-1]
	fileName, err := common.GetNVMDir("node", baseFileName)

	if err != nil {
		return
	}

	log.Println("filename", fileName)

	if err := utils.CheckSha(baseFileName, node_body, sha_body); err != nil {
		return err
	}

	dir := filepath.Dir(fileName)
	os.MkdirAll(dir, 0755)

	// TODO: stream
	if err := utils.SaveToFile(node_body, fileName); err != nil {
		return err
	}

	defer os.Remove(fileName)

	// TODO: join ungzip and untar

	unzipped, err := utils.UnGzip(fileName, dir)

	if err != nil {
		return
	}

	log.Println("UnGzip", unzipped)

	defer os.Remove(unzipped)

	if err := utils.Untar(unzipped, dir); err != nil {
		return err
	}

	if err != nil {
		return
	}

	// feels like I shouldn't care about this error
	targetDir, _ := common.GetNVMDir("node", string(version))
	sourceDir := strings.TrimSuffix(unzipped, ".tar")

	err = os.Rename(sourceDir, targetDir)

	defer os.RemoveAll(sourceDir)

	switch {
	case err == nil:
		break
	case errors.Is(err, os.ErrExist):
		// TODO: should we force overwrite?
		log.Println("file already existed")
	default:
		// returns some other error
		return fmt.Errorf("uncaught error: %w", err)
	}

	// create symlinks of bins
	// we've already established the function should work
	// at this point; ignoring error
	binTarget, _ := common.GetNVMDir("bin")
	// target all symlinks to nvmShim
	nvmShim := filepath.Join(binTarget, NVM_SHIM)
	// installed bins
	nodeBins := filepath.Join(targetDir, "bin")
	bins, err := os.ReadDir(nodeBins)

	if err != nil {
		return
	}

	for _, bin := range bins {
		source := filepath.Join(binTarget, bin.Name())

		err = syscall.Symlink(nvmShim, source)

		if err != nil && !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("cannot make symlink for %s = %w", bin.Name(), err)
		}
	}

	fmt.Println("Successfully installed", version)

	if !*installNoUse {
		return Use([]string{"--no-install", string(version)})
	}

	return nil
}
