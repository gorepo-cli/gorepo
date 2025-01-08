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

var CommandRegistry = []func(*config.Dependencies) *cli.Command{
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

func RegisterCommands(dependencies *config.Dependencies) []*cli.Command {
	var commands []*cli.Command
	for _, registerFunc := range CommandRegistry {
		commands = append(commands, registerFunc(dependencies))
	}
	return commands
}

var ModuleCommandRegistry = []func(string, *config.Dependencies) *cli.Command{
	exec.
		RegisterModuleCommand,
	fmt.
		RegisterModuleCommand,
	vet.
		RegisterModuleCommand,
}

func RegisterModuleCommands(moduleName string, dependencies *config.Dependencies) []*cli.Command {
	var commands []*cli.Command
	for _, registerFunc := range ModuleCommandRegistry {
		commands = append(commands, registerFunc(moduleName, dependencies))
	}
	return commands
}
