package commands

import (
	"flag"
	"fmt"

	common "github.com/bozdoz/nevermind/nvm-common"
)

// flag set for [commands.Which]
var whichCmd = flag.NewFlagSet("which", flag.ContinueOnError)

func init() {
	registerCommand(command{
		FlagSet: whichCmd,
		aliases: []string{"where"},
		help:    "get the path to the node executable for a given version",
		Handler: whichHandler,
	})
}

// get the path to the node executable for a given version
func whichHandler(cmd string, args []string) (err error) {
	whichCmd.Parse(args)
	args = whichCmd.Args()

	var version common.Version

	if len(args) < 1 {
		version, err = common.GetCurrentVersion()

		if err != nil {
			err = fmt.Errorf("%s command requires a single argument for version", cmd)
		}
	} else {
		version, err = common.GetVersion(args[0])
	}

	if err != nil {
		return
	}

	bin, err := common.GetNodeBin(version, "node")

	if err != nil {
		return fmt.Errorf("executable not found for version %s", version)
	}

	fmt.Println(bin)

	return
}
