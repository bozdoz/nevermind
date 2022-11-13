package commands

import (
	"errors"
	"flag"
	"fmt"
	"log"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/utils"
)

// flag set for [commands.Use]
var UseCmd = flag.NewFlagSet("use", flag.ContinueOnError)

var useHelp = UseCmd.Bool("help", false, `prints help text`)

// TODO: should the default be an env var?
var useNoInstall = UseCmd.Bool("no-install", false, "disable installing if nvm use cannot find version")

func init() {
	UseCmd.Usage = func() {
		utils.PrintTabs("\tuse\tuse a version of node")
	}
}

// use a version of node (updates config with desired version)
func Use(args []string) (err error) {
	UseCmd.Parse(args)
	args = UseCmd.Args()

	if *useHelp {
		// TODO: maybe verbose help message here
		UseCmd.Usage()
		utils.FlushTabs()
		return
	}

	if len(args) == 0 {
		return errors.New("use did not get arguments")
	}

	version, err := common.GetVersion(args[0])

	if err != nil {
		return
	}

	_, err = common.GetNodeBin(version, "node")

	if err != nil {
		msg := fmt.Sprintf("version not found: %s", version)

		if *useNoInstall {
			return fmt.Errorf(msg)
		}
		// try installing
		fmt.Println(msg)
		fmt.Println("installing...")

		err = Install([]string{"--use=false", string(version)})

		if err != nil {
			return fmt.Errorf("%s, and could not install - %w", msg, err)
		}
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

	fmt.Printf("You are now using node v%s\n", version)

	return nil
}
