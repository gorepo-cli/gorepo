package execute

import (
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/commands"
	"gorepo-cli/internal/commands/exec"
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

	// register commands
	registeredCommands := commands.RegisterCommands(dependencies)

	var commandNames []string
	for _, command := range registeredCommands {
		commandNames = append(commandNames, command.Name)
	}
	dependencies.Config.PushForbiddenNames(commandNames)

	app := &cli.App{
		Name:     "gorepo",
		Usage:    "A cli tool to manage Go monorepos",
		Commands: registeredCommands,
		CommandNotFound: func(c *cli.Context, command string) {
			l.FatalLn("command '" + command + "' not found")
		},
		Flags: flags.GlobalGroup,
	}

	if exists := cfg.RootConfigExists(); exists == true {
		// register root tasks as commands
		rootScripts, err := dependencies.Config.GetRootConfig()
		if err != nil {
			return err
		}
		if rootScripts.Tasks != nil && len(rootScripts.Tasks) > 0 {
			for taskName, _ := range rootScripts.Tasks {
				taskName := taskName
				app.Commands = append(app.Commands, exec.RegisterRootTaskCommand(taskName, dependencies))
			}
		}

		modules, err := cfg.GetModules([]string{"all"}, nil)
		if err != nil {
			return err
		}
		for _, module := range modules {
			module := module
			// register modules as commands
			registeredModuleCommands := commands.RegisterModuleCommands(module.Name, dependencies)

			// register module tasks as commands
			if module.Tasks != nil && len(module.Tasks) > 0 {
				for taskName, _ := range module.Tasks {
					taskName := taskName
					registeredModuleCommands = append(registeredModuleCommands, exec.RegisterModuleTaskCommand(module.Name, taskName, dependencies))
				}
			}

			app.Commands = append(app.Commands, &cli.Command{
				Name:        module.Name,
				Hidden:      true,
				Subcommands: registeredModuleCommands,
			})
		}
	}

	return app.Run(os.Args)
}
