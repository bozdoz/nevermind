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
	"time"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/utils"
)

const SHASUMS = "SHASUMS256.txt"

type node_download int

const (
	DOWNLOAD_NODE node_download = iota
	DOWNLOAD_SHASUM
)

// makes the enum exclusive
type exclusive_node_file interface {
	X()
}

func (d node_download) X() {}

// TODO: document env variables that we use
var BASE_URL = os.Getenv("NVM_NODEJS_ORG_MIRROR")
var InstallCmd = flag.NewFlagSet("install", flag.ContinueOnError)

func init() {
	if BASE_URL == "" {
		BASE_URL = "https://nodejs.org/dist"
	}

	InstallCmd.Usage = func() {
		utils.PrintTabs("\tinstall\tinstall a version of node")
		utils.PrintTabs("\t\talias: i")
	}
}

// this can be cached
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

func GetDownloadUrl(v string, d exclusive_node_file) string {
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

func Download(url string, ch chan []byte) (err error) {
	s := time.Now()
	res, err := http.Get(url)
	if err != nil {
		return
	}

	body, err := io.ReadAll(res.Body)

	res.Body.Close()

	log.Println("status", res.Status)
	log.Println("headers", res.Header)

	if res.StatusCode != 200 {
		// TODO: close channel?
		return fmt.Errorf("request of %s failed with status code: %d and\nbody: %s", url, res.StatusCode, body)
	}

	if err != nil {
		return
	}

	log.Printf("downloaded %s (%s)", url, time.Since(s))
	ch <- body

	return nil
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

func Install(args []string) (err error) {
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

	// TODO: allow cancelling the downloads via SIGINT
	go Download(node_url, node_chan)
	go Download(sha_url, sha_chan)

	node_body := <-node_chan
	sha_body := <-sha_chan

	segments := strings.Split(node_url, "/")
	baseFileName := segments[len(segments)-1]
	fileName := common.GetNVMDir("node", baseFileName)

	log.Println("filename", fileName)

	if err := CheckSha(baseFileName, node_body, sha_body); err != nil {
		return err
	}

	dir := filepath.Dir(fileName)
	os.MkdirAll(dir, 0755)

	if err := utils.SaveToFile(node_body, fileName); err != nil {
		return err
	}

	defer os.Remove(fileName)

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

	targetDir := common.GetNVMDir("node", version)
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
		return
	}

	log.Println("Success! Extracted node to", targetDir)

	return nil
}
