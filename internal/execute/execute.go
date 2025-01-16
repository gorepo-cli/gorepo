package execute

import (
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/commands"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"gorepo-cli/pkg"
	"gorepo-cli/pkg/logger"
	"os"
)

func Execute(l logger.Methods) (err error) {
	effects := pkg.NewEffects(l)
	cfg, err := config.NewConfig(effects)
	if err != nil {
		return err
	}
	dependencies := config.NewDependencies(effects, cfg)

	var cliCommands []*cli.Command

	cliCommands, err = commands.RegisterCommands(cliCommands, commands.RegisterRootCommands, commands.RegisterModuleCommands, dependencies)
	if err != nil {
		return err
	}

	var commandNames []string
	for _, command := range cliCommands {
		commandNames = append(commandNames, command.Name)
	}
	dependencies.Config.PushForbiddenNames(commandNames)

	app := &cli.App{
		Name:     "gorepo",
		Usage:    "A CLI tool to manage Go monorepos",
		Commands: cliCommands,
		CommandNotFound: func(c *cli.Context, command string) {
			l.FatalLn("command '" + command + "' not found")
		},
		Flags: flags.GlobalGroup,
	}

	return app.Run(os.Args)
}
