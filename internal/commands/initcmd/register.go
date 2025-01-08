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
	description = `Initialize a new monorepo at the working directory.

This command creates two primary files:
- 'work.toml' at the work directory
- 'go.work' file if the strategy is set as 'workspace' and one does not exist yet. This runs 'go work init' behind the hood`
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
