// nvm-shim intercepts requests to node, npm, and npx
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

// the main executable
const NODE = "node"

/*
LOCAL DEVELOPMENT use only
options are anything accessible in the node bin dir:
"node", "npm", "npx", "corepack"

set NVM_BIN=node at runtime
*/
var NVM_BIN = os.Getenv("NVM_BIN")

func init() {
	common.Debugger()
}

func fail(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func main() {
	bin, args := os.Args[0], os.Args[1:]

	bin = filepath.Base(bin)

	log.Println("NVM_BIN", NVM_BIN)
	log.Println("bin:", bin, "args:", args)

	// let's blow out any possible relative path exploits
	if strings.ContainsAny(bin, "./\\") {
		fail(fmt.Sprintf("invalid bin: %s", bin))
	}

	// might be able to extract GetCurrent from this
	config, err := common.GetConfig()

	if err != nil {
		fail("no node version installed. Did you run `nvm install`?")
	}

	version := config.Current

	if version == "" {
		fail("nvm encountered an error: version is empty! did you install a node version with `nvm install`?")
	}

	bin, err = common.GetNodeBin(version, bin)

	if err != nil {
		// check if running locally with NVM_BIN env var
		if NVM_BIN != "" {
			bin, err = common.GetNodeBin(version, NVM_BIN)
		}

		if err != nil {
			fail(fmt.Sprintf("failed to get bin from node version: %s, v%s — %s", bin, version, err))
		}
	}

	log.Println("running", bin)

	cmd := exec.Command(bin, args...)

	cmd.Stdin = os.Stdin

	if len(args) == 0 {
		// when args == 0, because I get:
		// exec: Stdout already set
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.Run()
	} else {
		out, err := cmd.Output()

		if err != nil {
			fail(err.Error())
		}

		if len(out) != 0 {
			fmt.Println(string(out))
		}
	}
}
