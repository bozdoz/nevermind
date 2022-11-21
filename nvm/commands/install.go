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
const install = "install"

var installCmd = flag.NewFlagSet(install, flag.ContinueOnError)
var installHelp = installCmd.Bool("help", false, "prints help text")
var installNoUse = installCmd.Bool("no-use", false, "disable calling nvm use <version> after install")
var installForce = installCmd.Bool("force", false, "force install")

func init() {
	registerCommand(command{
		FlagSet: installCmd,
		aliases: []string{"i"},
		help:    "install a version of node",
		Handler: installHandler,
	})
}

// install a given node version
// TODO? nvm install -h outputs help text twice
// TODO: break up into smaller functions
func installHandler(_ string, args []string) (err error) {
	// TODO: make parse work with flags after args
	installCmd.Parse(args)
	args = installCmd.Args()

	log.Println("parsed args", args)

	if *installHelp {
		// TODO: maybe verbose help message here
		installCmd.Usage()
		utils.FlushTabs()
		return
	}

	var version common.Version

	switch {
	case len(args) == 0:
		version, err = utils.ReadNvmrc()

		switch {
		case errors.Is(err, utils.ErrNoNvmRc):
			err = errors.New("install did not get arguments")
		case err != nil:
			err = fmt.Errorf("failed to get .nvmrc version: %w", err)
		}
	case args[0] == "lts":
		version, err = utils.GetLTS()
	default:
		version, err = common.GetVersion(args[0])
	}

	if err != nil {
		return
	}

	if !version.IsSpecific() {
		// look up latest within range
		version, err = utils.GetLatestFromVersion(version)
	}

	// check if version already installed, unless --force is present
	if !*installForce {
		// should be able to ignore error here
		installed, _ := utils.GetInstalledVersions()

		for _, ver := range installed {
			if ver == version {
				fmt.Printf("version %s already installed!\n", version)
				return
			}
		}
	}

	log.Println("installing version", version)

	node_url := utils.GetDownloadUrl(version, utils.DOWNLOAD_NODE)
	sha_url := utils.GetDownloadUrl(version, utils.DOWNLOAD_SHASUM)
	segments := strings.Split(node_url, "/")
	nodeFileName := segments[len(segments)-1]

	log.Println("node url", node_url)

	// AFAIK this is the best way to do async http requests
	node_chan := make(chan []byte)
	sha_chan := make(chan string)
	err_chan := make(chan error)

	fmt.Printf("Downloading node v%s\n", version)

	// TODO? allow cancelling the downloads via SIGINT
	go func() {
		node_download, err := utils.DownloadNode(node_url)

		if err != nil {
			err_chan <- err
		} else {
			node_chan <- node_download
		}
	}()

	go func() {
		sha, err := utils.FetchSha(sha_url, nodeFileName)

		if err != nil {
			err_chan <- err
		} else {
			sha_chan <- sha
		}
	}()

	var node_body []byte
	var sha_body string

	// waits on node request
	select {
	case node_body = <-node_chan:
	case err = <-err_chan:
		return
	}

	// waits on sha request
	select {
	case sha_body = <-sha_chan:
	case err = <-err_chan:
		return
	}

	err = utils.CheckSha(node_body, sha_body)

	if err != nil {
		return
	}

	fileName, err := common.GetNVMDir("node", nodeFileName)

	if err != nil {
		return
	}

	log.Println("filename", fileName)

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
		// TODO: force overwrite if "--force"
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

	fmt.Printf("Successfully installed v%s\n", version)

	if !*installNoUse {
		return useHandler(use, []string{"--no-install", string(version)})
	}

	return nil
}
