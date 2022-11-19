package commands

import (
	"flag"
	"fmt"
	"os"

	common "github.com/bozdoz/nevermind/nvm-common"
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
	dir, err := common.GetNVMDir("node")

	if err != nil {
		return
	}

	files, err := os.ReadDir(dir)

	if err != nil {
		return
	}

	dirs := make([]string, 0, len(files))

	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}

	fmt.Println(dirs)

	config, _ := common.GetConfig()

	if config.Current != "" {
		fmt.Println("Current", config.Current)
	}

	return
}
