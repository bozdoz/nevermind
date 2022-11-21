package commands

import (
	"flag"
	"fmt"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/utils"
)

// flag set for [commands.List]
const list = "list"

var listCmd = flag.NewFlagSet(list, flag.ContinueOnError)

func init() {
	registerCommand(command{
		FlagSet: listCmd,
		aliases: []string{"ls"},
		help:    "list installed node versions",
		Handler: listHandler,
	})
}

// list installed node versions
func listHandler(_ string, _ []string) (err error) {
	versions, err := utils.GetInstalledVersions()

	if err != nil {
		return
	}

	config, err := common.GetConfig()

	if err != nil {
		return
	}

	// TODO: should this be vertical list?
	fmt.Println(versions)

	if config.Current != "" {
		fmt.Println("Current", config.Current)
	}

	return
}
