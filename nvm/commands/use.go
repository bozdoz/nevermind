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

var useCmd = flag.NewFlagSet("use", flag.ContinueOnError)
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
	} else {
		// check if this specific version is installed
		_, err = common.GetNodeBin(version, "node")

		if err != nil {
			msg := fmt.Sprintf("version not found: %s", version)

			// --no-install
			if *useNoInstall {
				return fmt.Errorf(msg)
			}

			// try installing
			fmt.Println(msg)
			fmt.Println("installing...")

			_, err = install(version, installOptions{})

			if err == nil {
				// successful install also called `use`
				return
			} else {
				return fmt.Errorf("%s, and could not install - %w", msg, err)
			}
		}
	}

	return use(version)
}

// updates the current version in the config
//
// version here should be a specific version
func use(version common.Version) (err error) {
	// seems silly to verify version specificity again
	if !version.IsSpecific() {
		return fmt.Errorf("version is not specific, and cannot be used: %s", version)
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

		// make sure we have the binaries in our nvm PATH
		err = common.SyncSymlinks(version)

		if err != nil {
			return
		}
	}

	fmt.Printf("You are now using node v%s\n", version)

	return
}
