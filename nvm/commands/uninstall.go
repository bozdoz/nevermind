package commands

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/utils"
)

// flag set for [commands.Uninstall]
var UninstallCmd = flag.NewFlagSet("uninstall", flag.ContinueOnError)

func init() {
	UninstallCmd.Usage = func() {
		utils.PrintTabs("\tuninstall\tuninstall a specific node version")
	}
}

// uninstall a specific version of node
func Uninstall(args []string) (err error) {
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
