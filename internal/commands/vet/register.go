package vet

import (
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
)

var (
	name        = "vet"
	usage       = "Run go vet, break with --ci (module syntax compatible)"
	usageText   = "gorepo [global_options] [module_name] vet [command_options]"
	description = `Compatible with module syntax.

This command runs vet in all targeted modules.
It breaks without formating the files if you pass --ci.`
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
			return vet(dependencies, cmdFlags, globalFlags)
		},
		Flags: append(flags.ExecutionGroup, flags.Ci),
	}
}

func RegisterModuleCommand(moduleName string, dependencies *config.Dependencies) *cli.Command {
	return &cli.Command{
		Name:        name,
		Usage:       usage,
		UsageText:   usageText,
		Description: description,
		Action: func(c *cli.Context) error {
			globalFlags := flags.ExtractGlobalFlags(c)
			cmdFlags := flags.ExtractCommandFlags(c)
			cmdFlags.Target = moduleName
			return vet(dependencies, cmdFlags, globalFlags)
		},
		Flags: []cli.Flag{flags.Ci},
	}
}
