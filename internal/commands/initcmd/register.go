package initcmd

import (
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
)

var (
	name        = "init"
	usage       = "Initialize a new monorepo"
	usageText   = "gorepo [--global_options] init [--command_options] [monorepo_name]"
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
			return initCmd(dependencies, cmdFlags, globalFlags, c.Args().Get(0))
		},
	}
}
