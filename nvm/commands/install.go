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
	"time"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/utils"
)

var installCmd = flag.NewFlagSet("install", flag.ContinueOnError)
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

	_, err = install(version, installOptions{
		force: *installForce,
		noUse: *installNoUse,
	})

	return err
}

type installOptions struct {
	force, noUse bool
}

func install(version common.Version, options installOptions) (ret common.Version, err error) {
	if !version.IsSpecific() {
		// look up latest within range
		version, err = utils.GetLatestFromVersion(version)
	}

	// check if version already installed, unless --force is present
	if !options.force {
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

	nodeDir, err := common.GetNVMDir("node")

	if err != nil {
		return
	}

	os.MkdirAll(nodeDir, 0755)

	s := time.Now()

	err = utils.UnArchiveBytes(node_body, nodeDir)

	if err != nil {
		return
	}

	log.Println("time to save files:", time.Since(s))

	targetDir := filepath.Join(nodeDir, string(version))
	sourceDir := filepath.Join(nodeDir, strings.TrimSuffix(nodeFileName, ".tar.gz"))

	err = os.Rename(sourceDir, targetDir)

	switch {
	case err == nil:
		defer os.RemoveAll(sourceDir)
	case errors.Is(err, os.ErrExist):
		// if it exists, remove it and replace it
		// TODO test this
		os.RemoveAll(targetDir)
		err = os.Rename(sourceDir, targetDir)

		if err != nil {
			// dang
			return
		}
	default:
		// returns some other error
		return ret, fmt.Errorf("uncaught error: %w", err)
	}

	success := func() {
		fmt.Printf("Successfully installed v%s\n", version)
	}

	if !options.noUse {
		// install is done
		success()
		// `use` syncs symlinks
		err = use(version)
	} else {
		success()
	}

	return
}
