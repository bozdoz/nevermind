/*
nevermind. a node version manager

# Environment Variables

DEBUG:
  - set DEBUG=* to enable ALL debugging
  - set DEBUG=1 to enable go debugging

NVM_NODEJS_ORG_MIRROR:
  - set an alternate proxy for https://nodejs.org/dist
*/
package main

//go:generate go build -o $HOME/.nevermind/bin/nvm ./

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/commands"
	"github.com/bozdoz/nevermind/nvm/utils"
)

// nvm version
const VERSION = "v0.1.5"

var help *flag.FlagSet

var vFlag = flag.Bool("v", false, "print version of nvm")

func init() {
	common.Debugger()

	// uses tabwriter to align tabs
	flag.CommandLine.SetOutput(utils.Writer)

	flag.Usage = func() {
		utils.PrintTabs("nevermind. a node version manager.")
		utils.PrintTabs("")
		utils.PrintTabs("USAGE:")
		utils.PrintTabs("\tnvm install <version>")
		utils.PrintTabs("\tnvm use <version>")
		utils.PrintTabs("")
		utils.PrintTabs("FLAGS:")
		flag.PrintDefaults()
		utils.PrintTabs("")
		utils.PrintTabs("SUBCOMMANDS:")

		commands.Commands.Usage()

		utils.FlushTabs()
	}
}

func fail(message string) {
	fmt.Println("Error:", message)

	if help == nil {
		flag.Usage()
	} else {
		help.Usage()
	}

	fmt.Println("")

	utils.FlushTabs()
	os.Exit(1)
}

// TODO: a command to show executables from NVM_DIR/bin
func main() {
	flag.Parse()

	args := flag.Args()

	log.Println("args", args)

	if len(args) == 0 {
		// -v passed
		if *vFlag {
			fmt.Println(VERSION)
		} else {
			fail("Please specify a subcommand\n")
		}
		return
	}

	cmd, args := args[0], args[1:]

	var err error

	// TODO: command to show config
	found, ok := commands.Commands[cmd]

	if ok {
		help = found.FlagSet

		err = found.Handler(cmd, args)
	} else {
		err = fmt.Errorf("subcommand does not exist: %s", cmd)
	}

	if err != nil {
		fail(err.Error())
	}
}
