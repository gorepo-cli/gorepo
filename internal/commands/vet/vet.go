package vet

import (
	"errors"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"path/filepath"
	"strings"
)

func vet(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags) error {
	if err := dependencies.Config.BreakIfRootConfigDoesNotExist(); err != nil {
		return err
	}

	verbose := globalFlags.Verbose
	if verbose {
		dependencies.Effects.Logger.VerboseLn("verbose mode enabled")
	}

	targets := strings.Split(cmdFlags.Target, ",")
	if verbose {
		dependencies.Effects.Logger.VerboseLn("value for flag target:       " + strings.Join(targets, ","))
	}

	exclude := strings.Split(cmdFlags.Exclude, ",")
	if verbose {
		dependencies.Effects.Logger.VerboseLn("value for flag exclude:      " + strings.Join(exclude, ","))
	}

	// executing the script at the root means executing it in all modules
	if targets[0] == "root" {
		targets = []string{"all"}
	}

	modules, err := dependencies.Config.GetModules(targets, exclude)
	if err != nil {
		return err
	}

	script := "go vet . || exit 1"

	for _, module := range modules {
		path := filepath.Join(dependencies.Config.Runtime.ROOT, module.RelativePath)
		if module.Language != "go" {
			if err := dependencies.Effects.Executor.Bash(path, script); err != nil {
				return errors.New("error: vet --ci failed in module " + module.Name)
			}
		} else {
			hasVet := false
			for name, _ := range module.Scripts {
				if name == "vet" {
					hasVet = true
				}
			}
			if hasVet {
				dependencies.Effects.Logger.InfoLn("module " + module.Name + " in not a go module but it implements vet, executing")
				for i, _ := range module.Scripts["vet"] {
					if err := dependencies.Effects.Executor.Bash(path, module.Scripts["fmt-ci"][i]); err != nil {
						return errors.New("error: fmt-ci failed in module " + module.Name)
					}
				}
			} else {
				dependencies.Effects.Logger.WarningLn("module " + module.Name + " in not a go module and it does not implement vet, skipping")
			}
		}
	}

	return nil
}
