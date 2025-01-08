package add

import (
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
)

var (
	name        = "add"
	usage       = "Add a module"
	usageText   = "gorepo [global_options] add [command_options] <module_name>"
	description = ""
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
