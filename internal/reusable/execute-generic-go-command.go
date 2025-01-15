package reusable

import (
	"errors"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"path/filepath"
	"strings"
)

func ExecuteGenericGoCommand(
	dependencies *config.Dependencies,
	cmdFlags *flags.CommandFlags,
	globalFlags *flags.GlobalFlags,
	commandName string,
	script string,
) error {
	if globalFlags.Verbose {
		dependencies.Effects.Logger.VerboseLn("the value of --target is '" + cmdFlags.Target + "'")
		dependencies.Effects.Logger.VerboseLn("the value of --exclude is '" + cmdFlags.Exclude + "'")
	}

	targets := strings.Split(cmdFlags.Target, ",")
	exclude := strings.Split(cmdFlags.Exclude, ",")

	if targets[0] == "root" {
		rootConfig, err := dependencies.Config.GetRootConfig()
		if err != nil {
			return err
		}
		if rootConfig.Tasks[commandName] != nil && len(rootConfig.Tasks[commandName]) > 0 {
			for i, _ := range rootConfig.Tasks[commandName] {
				if err := dependencies.Effects.Executor.Bash(dependencies.Config.Runtime.ROOT, rootConfig.Tasks[commandName][i]); err != nil {
					return errors.New("task " + commandName + " failed at root ")
				}
			}
			return nil
		}
		// if there is no task at root, target all modules
		targets = []string{"all"}
	}

	modules, err := dependencies.Config.GetModules(targets, exclude)
	if err != nil {
		return err
	}

	for _, moduleConfig := range modules {
		path := filepath.Join(dependencies.Config.Runtime.ROOT, moduleConfig.PathFromRoot)

		if moduleConfig.Tasks[commandName] != nil && len(moduleConfig.Tasks[commandName]) > 0 {
			for i, _ := range moduleConfig.Tasks[commandName] {
				if err := dependencies.Effects.Executor.Bash(dependencies.Config.Runtime.ROOT, moduleConfig.Tasks[commandName][i]); err != nil {
					return errors.New("task " + commandName + " failed at module " + moduleConfig.Name)
				}
			}
			continue
		}

		if moduleConfig.Language == "go" {
			if err := dependencies.Effects.Executor.Go(path, commandName); err != nil {
				return errors.New("command " + commandName + " failed at module " + moduleConfig.Name)
			}
			continue
		}
	}

	return nil
}
