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
const VERSION = "v0.1.2"

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
		commands.InstallCmd.Usage()
		commands.UninstallCmd.Usage()
		commands.UseCmd.Usage()
		commands.ListCmd.Usage()
		commands.WhichCmd.Usage()
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

func main() {
	flag.Parse()

	args := flag.Args()

	log.Println("args", args)
	log.Println("vFlag", *vFlag)

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

	// TODO: nvm versioning, show config
	// TODO: maybe we should time the whole command
	switch cmd {
	case "install", "i":
		help = commands.InstallCmd
		err = commands.Install(args)
	case "uninstall":
		help = commands.UninstallCmd
		err = commands.Uninstall(args)
	case "use":
		help = commands.UseCmd
		err = commands.Use(args)
	case "list", "ls":
		help = commands.ListCmd
		err = commands.List(cmd, args)
	case "which", "where":
		help = commands.WhichCmd
		err = commands.Which(cmd, args)
	default:
		err = fmt.Errorf("no command defined: %s", cmd)
	}

	if err != nil {
		fail(err.Error())
	}
}
