package commands

import (
	"errors"
	"flag"
	"fmt"
	"log"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/utils"
)

var UseCmd = flag.NewFlagSet("use", flag.ContinueOnError)

func init() {
	UseCmd.Usage = func() {
		utils.PrintTabs("\tuse\tuse a version of node")
	}
}

func Use(args []string) (err error) {
	if len(args) == 0 {
		return errors.New("use did not get arguments")
	}

	version, err := common.GetVersion(args[0])

	if err != nil {
		return
	}

	_, err = common.GetNodeBin(version, "node")

	if err != nil {
		// TODO: conditionally run install
		return fmt.Errorf("version not found: %s", version)
	}

	config, err := common.GetConfig()

	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	if config.Current == version {
		log.Println("version already set", config)
	} else {
		config.Current = version

		err = common.SetConfig(config)

		if err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	}

	return nil
}
