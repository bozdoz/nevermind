// nvm-shim intercepts requests to node, npm, npx, and any other
// binaries installed by npm
package main

//go:generate go build -o $HOME/.nevermind/bin/nvm-shim ./

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	common "github.com/bozdoz/nevermind/nvm-common"
)

func init() {
	common.Debugger()
}

func fail(message string) {
	log.Println("nvm-shim failed")
	fmt.Println(message)
	os.Exit(1)
}

func main() {
	bin, args := os.Args[0], os.Args[1:]

	bin = filepath.Base(bin)

	log.Println("bin:", bin, "args:", args)

	// let's blow out any possible relative path exploits
	if strings.ContainsAny(bin, "$~./\\") {
		fail(fmt.Sprintf("invalid bin: %s", bin))
	}

	version, err := common.GetCurrentVersion()

	if err != nil {
		fail("no node version installed. Did you run `nvm install`?")
	}

	if version == "" {
		fail("nvm encountered an error: version is empty! did you install a node version with `nvm install`?")
	}

	absBin, err := common.GetNodeBin(version, bin)

	if err != nil {
		fail(fmt.Sprintf("failed to get bin from node version: %s, v%s — %s", bin, version, err))
	}

	log.Println("running", absBin)

	cmd := exec.Command(absBin, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	// symlink newly installed binaries
	if bin == "npm" && isGlobalInstall(args) {
		// check for newly installed binaries
		log.Println("creating new symlinks for binaries")
		common.SyncSymlinks(version)
	}

	// give npm --help what it wants
	if err != nil {
		os.Exit(1)
	}
}

var globalFlags = map[string]bool{
	"-g":                true,
	"--global":          true,
	"--location=global": true,
}

// from `npm i --help`
var installAliases = map[string]bool{
	"add":     true,
	"i":       true,
	"in":      true,
	"ins":     true,
	"inst":    true,
	"insta":   true,
	"instal":  true,
	"install": true,
	"isnt":    true,
	"isnta":   true,
	"isntal":  true,
	"isntall": true,
}

// checks if we're installing a global binary
func isGlobalInstall(args []string) bool {
	var isGlobal bool
	var isInstall bool

	for i, arg := range args {
		if globalFlags[arg] {
			isGlobal = true
		}
		if arg == "--location" && args[i+1] == "global" {
			isGlobal = true
		}
		if installAliases[arg] {
			isInstall = true
		}
	}

	return isGlobal && isInstall
}
