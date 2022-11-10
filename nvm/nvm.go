package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/commands"
	"github.com/bozdoz/nevermind/nvm/utils"
)

type Env struct {
	Help *flag.FlagSet
}

var env = &Env{}

func init() {
	common.Debugger()

	flag.Usage = func() {
		utils.PrintTabs("nevermind. a node version manager.")
		utils.PrintTabs("")
		utils.PrintTabs("USAGE:")
		utils.PrintTabs("\tnvm install <version>")
		utils.PrintTabs("\tnvm use <version>")
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
	fmt.Println("")
	fmt.Println(message)

	if env.Help == nil {
		flag.Usage()
	} else {
		env.Help.Usage()
	}

	fmt.Println("")

	utils.FlushTabs()
	os.Exit(1)
}

func main() {
	flag.Parse()

	args := flag.Args()

	log.Println("running with args", args)

	if len(args) == 0 {
		fmt.Println("Please specify a subcommand.")
		fmt.Println("")
		flag.Usage()
		os.Exit(1)
	}

	cmd, args := args[0], args[1:]

	var err error

	// TODO: nvm versioning, show config
	// TODO: maybe we should time the whole command
	switch cmd {
	case "install", "i":
		env.Help = commands.InstallCmd
		err = commands.Install(args)
	case "uninstall":
		env.Help = commands.UninstallCmd
		err = commands.Uninstall(args)
	case "use":
		env.Help = commands.UseCmd
		err = commands.Use(args)
	case "list", "ls":
		env.Help = commands.ListCmd
		err = commands.List(cmd, args)
	case "which", "where":
		env.Help = commands.WhichCmd
		err = commands.Which(cmd, args)
	default:
		err = fmt.Errorf("no command defined: %s", cmd)
	}

	if err != nil {
		fail(err.Error())
	}
}
