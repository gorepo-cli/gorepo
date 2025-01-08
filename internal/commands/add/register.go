package add

import (
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
)

var (
	name        = "add"
	usage       = "Add a module"
	usageText   = "gorepo [--global_options] add [--command_options] <module_name>"
	description = `Add a new module to the monorepo.

This command creates a new folder with 2 file, 'module.toml' and 'go.mod'. It also adds the module to the go workspace. You can pass a path ending with the module name.`
)

func RegisterCommand(dependencies *config.Dependencies) *cli.Command {
	return &cli.Command{
		Name:        name,
		Usage:       usage,
		UsageText:   usageText,
		Description: description,
		Action: func(c *cli.Context) error {
			globalFlags := flags.ExtractGlobalFlags(c)
			cmdFlags := flags.ExtractCommandFlags(c)
			relativePathAndNameInput := c.Args().Get(0)
			return add(dependencies, cmdFlags, globalFlags, relativePathAndNameInput)
		},
	}
}
