package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	common "github.com/bozdoz/nevermind/nvm-common"
)

func init() {
	common.Debugger()
}

func fail(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func main() {
	args := os.Args[1:]

	// might be able to extract GetCurrent from this
	config, err := common.GetConfig()

	if err != nil {
		fail(err.Error())
	}

	version := config.Current

	if version == "" {
		fail("version is empty! did you install a node version?")
	}

	node, err := common.GetNodeBin(version, "node")

	if err != nil {
		fail(fmt.Sprintf("failed to get node at version: %s. %s", version, err))
	}

	log.Println("running", node, args)

	cmd := exec.Command(node, args...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if len(args) == 0 {
		cmd.Run()
	} else {
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if len(out) != 0 {
			fmt.Println(string(out))
		}
	}
}
