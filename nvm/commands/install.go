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
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/utils"
)

// remote filename which indicates sha's for each downloadable file
const SHASUMS = "SHASUMS256.txt"

// default download location
const DEFAULT_BASE_URL = "https://nodejs.org/dist"

// basename of binary in NVM_DIR
const NVM_SHIM = "nvm-shim"

// Enum for downloading either node or shasums
type NodeDownload int

const (
	DOWNLOAD_NODE   NodeDownload = iota // downloads node
	DOWNLOAD_SHASUM                     // downloads SHASUMS
)

// makes the enum exclusive
type exclusive_node_file interface {
	X()
}

// ignore me
func (d NodeDownload) X() {}

// TODO: document env variables that we use

// NVM_NODEJS_ORG_MIRROR - base url for downloading node; default: "https://nodejs.org/dist"
var BASE_URL = os.Getenv("NVM_NODEJS_ORG_MIRROR")

// flag set for [commands.Install]
var InstallCmd = flag.NewFlagSet("install", flag.ContinueOnError)

var installHelp = InstallCmd.Bool("help", false, `prints help text`)
var installNoUse = InstallCmd.Bool("no-use", false, `disable calling nvm use <version> after install`)

func init() {
	// sets default for BASE_URL
	if BASE_URL == "" {
		BASE_URL = DEFAULT_BASE_URL
	}

	InstallCmd.Usage = func() {
		utils.PrintTabs("\tinstall, i\tinstall a version of node")
	}
}

// TODO: this can be cached, or run once
// used to determine download url in [commands.GetDownloadUrl]
func GetOsAndArch() (remote_os, remote_arch string) {
	remote_arch = runtime.GOARCH

	switch remote_arch {
	case "x86_64", "amd64":
		remote_arch = "x64"
	case "aarch64":
		remote_arch = "arm64"
	}

	remote_os = runtime.GOOS

	if remote_os == "windows" {
		remote_os = "win"
		if remote_arch == "amd64" {
			remote_arch = "x64"
		} else {
			remote_arch = "x86"
		}
	}

	return
}

/*
determine where to download node install files

`v` a verified/parsed version string

`d` is a [commands.NodeDownload]

used by [commands.Install]
*/
func GetDownloadUrl(v common.Version, d exclusive_node_file) string {
	remote_os, remote_arch := GetOsAndArch()
	ext := "tar.gz"

	// TODO: we don't actually test extracting zips on windows
	if remote_os == "win" {
		ext = "zip"
		panic("we haven't written any code that would extract a zip yet")
	}

	switch d {
	case DOWNLOAD_NODE:
		return fmt.Sprintf("%s/v%s/node-v%s-%s-%s.%s", BASE_URL, v, v, remote_os, remote_arch, ext)
	case DOWNLOAD_SHASUM:
		return fmt.Sprintf("%s/v%s/%s", BASE_URL, v, SHASUMS)
	default:
		return ""
	}
}

/*
actually downloads a url and passes the []bytes
to the channel;
if there is an error, we pass it to `err_ch`
used by [commands.Install]
*/
func Download(url string, ch chan []byte, err_ch chan error) {
	defer close(ch)
	s := time.Now()
	res, err := http.Get(url)

	if err != nil {
		err_ch <- err
		return
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		err_ch <- err
		return
	}

	res.Body.Close()

	log.Println("status", res.Status)
	log.Println("headers", res.Header)

	if res.StatusCode != 200 {
		err_ch <- fmt.Errorf("request of %s failed with status code: %d", url, res.StatusCode)
		return
	}

	log.Printf("downloaded %s (%s)", url, time.Since(s))

	ch <- body
}

func CheckSha(fileName string, node_body, sha_body []byte) error {
	h := sha256.New()
	h.Write(node_body)
	file_sha := fmt.Sprintf("%x", h.Sum(nil))

	log.Println("looking for", file_sha)

	shas := strings.Fields(string(sha_body))

	for i, line := range shas {
		if line == fileName {
			verified_sha := shas[i-1]
			log.Println("found", verified_sha)

			if verified_sha != file_sha {
				return fmt.Errorf("no SHA match for %s", fileName)
			} else {
				log.Println("SHA's match!")

				return nil
			}
		}
	}

	return fmt.Errorf("could not find filename (%s) in %s", fileName, SHASUMS)
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

	node_url := GetDownloadUrl(version, DOWNLOAD_NODE)
	sha_url := GetDownloadUrl(version, DOWNLOAD_SHASUM)
	log.Println("download node url", node_url)
	log.Println("download sha url", sha_url)

	node_chan := make(chan []byte)
	sha_chan := make(chan []byte)
	err_chan := make(chan error)

	// TODO: allow cancelling the downloads via SIGINT
	go Download(node_url, node_chan, err_chan)
	go Download(sha_url, sha_chan, err_chan)

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
	fileName := common.GetNVMDir("node", baseFileName)

	log.Println("filename", fileName)

	if err := CheckSha(baseFileName, node_body, sha_body); err != nil {
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

	targetDir := common.GetNVMDir("node", string(version))
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
	binTarget := common.GetNVMDir("bin")
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
