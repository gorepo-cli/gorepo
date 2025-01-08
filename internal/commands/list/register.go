package list

import (
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
)

var (
	name        = "list"
	usage       = "List modules"
	usageText   = "gorepo [--global_options] list [--command_options]"
	description = `List all modules of the monorepo. Formally a module is a folder with a module.toml file in it, regardless of the language it uses.`
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
			return list(dependencies, cmdFlags, globalFlags)
		},
	}
}
