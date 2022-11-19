package commands

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/bozdoz/nevermind/nvm/utils"
)

type commandList map[string]command

// prints sorted usage of all subcommands
func (list commandList) Usage() {
	// omit aliases
	printed := map[string]bool{}

	keys := make([]string, 0, len(list))
	for k := range list {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		cmd := list[k].FlagSet
		name := cmd.Name()
		_, ok := printed[name]
		if !ok {
			cmd.Usage()
			printed[name] = true
		}
	}
}

// map all commands used to easily call functions related to flags
var Commands = commandList{}

type command struct {
	FlagSet *flag.FlagSet
	aliases []string
	help    string
	Handler func(subcmd string, args []string) (err error)
}

// TODO add common --help in handler
func registerCommand(cmd command) {
	name := cmd.FlagSet.Name()

	// use cmd.help in Usage for name and all aliases
	cmd.FlagSet.Usage = func() {
		flags := strings.Join(append([]string{name}, cmd.aliases...), ", ")

		utils.PrintTabs(fmt.Sprintf("\t%s\t%s", flags, cmd.help))
	}

	Commands[name] = cmd

	for _, v := range cmd.aliases {
		Commands[v] = cmd
	}
}
