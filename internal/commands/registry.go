package commands

import (
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/commands/add"
	"gorepo-cli/internal/commands/check"
	"gorepo-cli/internal/commands/exec"
	"gorepo-cli/internal/commands/fmt"
	"gorepo-cli/internal/commands/initcmd"
	"gorepo-cli/internal/commands/list"
	"gorepo-cli/internal/commands/version"
	"gorepo-cli/internal/commands/vet"
	"gorepo-cli/internal/config"
)

var commandRegistry = []func(*config.Dependencies) *cli.Command{
	initcmd.
		RegisterCommand,
	add.
		RegisterCommand,
	list.
		RegisterCommand,
	exec.
		RegisterCommand,
	fmt.
		RegisterCommand,
	vet.
		RegisterCommand,
	check.
		RegisterCommand,
	version.
		RegisterCommand,
}

func RegisterRootCommands(dependencies *config.Dependencies) []*cli.Command {
	var commands []*cli.Command

	for _, registerFunc := range commandRegistry {
		commands = append(commands, registerFunc(dependencies))
	}

	return commands
}

var moduleCommandRegistry = []func(string, *config.Dependencies) *cli.Command{
	exec.
		RegisterModuleCommand,
	fmt.
		RegisterModuleCommand,
	vet.
		RegisterModuleCommand,
}

func RegisterModuleCommands(moduleName string, dependencies *config.Dependencies) []*cli.Command {
	var commands []*cli.Command

	for _, registerFunc := range moduleCommandRegistry {
		commands = append(commands, registerFunc(moduleName, dependencies))
	}

	return commands
}

func RegisterCommands(
	c []*cli.Command,
	registerRootCommands func(dependencies *config.Dependencies) []*cli.Command,
	registerModuleCommands func(moduleName string, dependencies *config.Dependencies) []*cli.Command,
	dependencies *config.Dependencies) ([]*cli.Command, error) {
	commands := make([]*cli.Command, 0)

	if c != nil {
		commands = c
	}

	var commandNames []string

	rootCommands := registerRootCommands(dependencies)

	for _, command := range rootCommands {
		commands = append(commands, command)
		commandNames = append(commandNames, command.Name)
	}

	if exists := dependencies.Config.RootConfigExists(); exists == true {
		modules, err := dependencies.Config.GetModules([]string{"all"}, nil)
		if err != nil {
			return nil, err
		}

		for _, module := range modules {
			module := module

			commandExists := false
			for _, c := range commandNames {
				if c == module.Name {
					commandExists = true
					break
				}
			}
			if commandExists {
				//dependencies.Effects.Logger.WarningLn("not registering module '" + module.Name + "' as a command because a command has this name already. This is not critical, but you will not be able to use syntactic sugar for this module")
				continue
			}

			moduleCommands := registerModuleCommands(module.Name, dependencies)

			if module.Tasks != nil && len(module.Tasks) > 0 {
				for taskName, _ := range module.Tasks {
					taskName := taskName
					commandExists := false
					for _, c := range moduleCommands {
						if c.Name == taskName {
							commandExists = true
							break
						}
					}
					if commandExists {
						//dependencies.Effects.Logger.WarningLn("not registering task '" + taskName + "' as a subcommand for module " + module.Name + " because a subcommand has this name already. This is not critical, but you will not be able to use syntactic sugar for this task")
						continue
					}
					moduleCommands = append(moduleCommands, exec.RegisterModuleTaskCommand(module.Name, taskName, dependencies))
				}
			}

			commandNames = append(commandNames, module.Name)
			commands = append(commands, &cli.Command{
				Name:        module.Name,
				Hidden:      true,
				Subcommands: moduleCommands,
			})
		}

		rootScripts, err := dependencies.Config.GetRootConfig()
		if err != nil {
			return nil, err
		}

		if rootScripts.Tasks != nil && len(rootScripts.Tasks) > 0 {
			for taskName, _ := range rootScripts.Tasks {
				taskName := taskName
				commandExists := false
				for _, c := range commandNames {
					if c == taskName {
						commandExists = true
						break
					}
				}
				if commandExists {
					//dependencies.Effects.Logger.WarningLn("not registering task '" + taskName + "' as a command because a command has this name already. This is not critical, but you will not be able to use syntactic sugar for this task")
					continue
				}
				commandNames = append(commandNames, taskName)
				commands = append(commands, exec.RegisterRootTaskCommand(taskName, dependencies))
			}
		}
	}

	return commands, nil
}
