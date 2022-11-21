package commands

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	common "github.com/bozdoz/nevermind/nvm-common"
	"github.com/bozdoz/nevermind/nvm/utils"
)

const use = "use"

var useCmd = flag.NewFlagSet(use, flag.ContinueOnError)
var useHelp = useCmd.Bool("help", false, `prints help text`)

// TODO: should the default be an env var?
var useNoInstall = useCmd.Bool("no-install", false, "disable installing if nvm use cannot find version")

func init() {
	registerCommand(command{
		FlagSet: useCmd,
		help:    "use a version of node",
		Handler: useHandler,
	})
}

// use a version of node (updates config with desired version)
func useHandler(_ string, args []string) (err error) {
	// TODO: maybe move to generic command registry
	useCmd.Parse(args)
	args = useCmd.Args()

	if *useHelp {
		// TODO: maybe verbose help message here
		useCmd.Usage()
		utils.FlushTabs()
		return
	}

	var version common.Version

	if len(args) == 0 {
		version, err = utils.ReadNvmrc()

		if err != nil {
			return errors.New("use did not get arguments")
		}
	} else {
		// version in cli arg
		version, err = common.GetVersion(args[0])
	}

	if err != nil {
		return
	}

	// if not specific, then check list to see if we have it
	if !version.IsSpecific() {
		installed, err := utils.GetInstalledVersions()

		if err == nil {
			// check list to see if we already have a match
			for _, ver := range installed {
				if strings.HasPrefix(string(ver), fmt.Sprintf("%s.", version)) {
					version = ver
				}
			}
		}
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

		// TODO: how to get installed version from non-specific
		err = installHandler(install, []string{"--no-use", string(version)})

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

	// TODO: this should be the actual specific version; not whatever was passed
	fmt.Printf("You are now using node v%s\n", version)

	return nil
}
