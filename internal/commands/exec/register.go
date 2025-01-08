package exec

import (
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
)

var (
	name        = "exec"
	usage       = "Execute a script"
	usageText   = "gorepo [global_options] [module_name] exec [command_options] <script_name>"
	description = `Compatible with module syntax.

Execute a script at the root of the monorepo, or in one, many or all modules. Scripts are declared in the files 'work.toml' and 'module.toml'.`
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
			return exec(dependencies, cmdFlags, globalFlags, c.Args().Get(0))
		},
		Flags: append(flags.ExecutionGroup, flags.AllowMissing),
	}
}

func RegisterModuleCommand(moduleName string, dependencies *config.Dependencies) *cli.Command {
	return &cli.Command{
		Name:   name,
		Hidden: true,
		Action: func(c *cli.Context) error {
			globalFlags := flags.ExtractGlobalFlags(c)
			cmdFlags := flags.ExtractCommandFlags(c)
			cmdFlags.Target = moduleName
			return exec(dependencies, cmdFlags, globalFlags, c.Args().Get(0))
		},
	}
}
