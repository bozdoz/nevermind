package commands

import (
	"flag"
	"fmt"
	"os"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/utils"
)

// flag set for [commands.List]
var ListCmd = flag.NewFlagSet("list", flag.ContinueOnError)

func init() {
	ListCmd.Usage = func() {
		utils.PrintTabs("\tlist, ls\tlist installed node versions")
	}
}

// list installed node versions
func List(cmd string, args []string) (err error) {
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
