package commands

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	common "github.com/bozdoz/nevermind/nvm-common"
)

const uninstall = "uninstall"

var uninstallCmd = flag.NewFlagSet(uninstall, flag.ContinueOnError)

func init() {
	registerCommand(command{
		FlagSet: uninstallCmd,
		help:    "uninstall a specific node version",
		Handler: uninstallHandler,
	})
}

// uninstall a specific version of node
func uninstallHandler(_ string, args []string) (err error) {
	uninstallCmd.Parse(args)
	args = uninstallCmd.Args()

	if len(args) < 1 {
		return errors.New("uninstall requires a single argument for version")
	}
	version, err := common.GetVersion(args[0])

	if err != nil {
		return
	}

	_, err = common.GetNodeBin(version, "node")

	if err != nil {
		return fmt.Errorf("executable not found for version %s", version)
	}

	// unset config
	config, err := common.GetConfig()

	if err == nil {
		// TODO: what should the fallback version be?
		// should Version type be able to be invalid?
		config.Current = ""
		err = common.SetConfig(config)

		if err != nil {
			log.Println("failed to set config: %w", err)
		}
	} else {
		log.Println("failed to get config: %w", err)
	}

	installDir, err := common.GetNVMDir("node", string(version))

	if err != nil {
		return
	}

	err = os.RemoveAll(installDir)

	if err != nil {
		return
	}

	fmt.Println("Successfully uninstalled", version)

	return nil
}
