package vet

import (
	"errors"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"path/filepath"
	"strings"
)

// todo: broken because of ci

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

	if targets[0] == "root" {
		return errors.New("running vet-ci from root is not supported")
	}

	modules, err := dependencies.Config.GetModules(targets, exclude)
	if err != nil {
		return err
	}

	script := "go vet . || exit 1"

	for _, module := range modules {
		path := filepath.Join(dependencies.Config.Runtime.ROOT, module.RelativePath)
		if err := dependencies.Effects.Executor.Bash(path, script); err != nil {
			return errors.New("error: vet --ci failed in module " + module.Name)
		}
	}

	return nil
}
