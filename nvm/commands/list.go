package commands

import (
	"flag"
	"fmt"
	"os"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/utils"
)

var ListCmd = flag.NewFlagSet("list", flag.ContinueOnError)

func init() {
	ListCmd.Usage = func() {
		utils.PrintTabs("\tlist\tlist installed node versions")
		utils.PrintTabs("\t\talias: ls")
	}
}

func List(cmd string, args []string) (err error) {
	dir := common.GetNVMDir("node")

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
