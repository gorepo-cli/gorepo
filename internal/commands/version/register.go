package version

import (
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
)

var (
	name        = "version"
	usage       = "Print version"
	usageText   = "gorepo [global_options] version [command_options]"
	description = "(module syntax compatible)"
)

func RegisterCommand(dependencies *config.Dependencies) *cli.Command {
	return &cli.Command{
		Name:        name,
		Usage:       usage,
		UsageText:   usageText,
		Description: description,
		Action: func(c *cli.Context) error {
			globalFlags := flags.ExtractGlobalFlags(c)
			return version(dependencies, nil, globalFlags)
		},
	}
}
