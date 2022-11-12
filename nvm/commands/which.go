package commands

import (
	"flag"
	"fmt"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/utils"
)

// flag set for [commands.Which]
var WhichCmd = flag.NewFlagSet("which", flag.ContinueOnError)

func init() {
	WhichCmd.Usage = func() {
		utils.PrintTabs("\twhich, where\tget the path to the node executable for a given version")
	}
}

// get the path to the node executable for a given version
func Which(cmd string, args []string) (err error) {
	if len(args) < 1 {
		return fmt.Errorf("%s command requires a single argument for version", cmd)
	}
	version, err := common.GetVersion(args[0])

	if err != nil {
		return
	}

	bin, err := common.GetNodeBin(version, "node")

	if err != nil {
		return fmt.Errorf("executable not found for version %s", version)
	}

	fmt.Println(bin)

	return nil
}
