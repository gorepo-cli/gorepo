package fmt

import (
	"errors"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"path/filepath"
	"strings"
)

// todo: broken because of ci

func fmt(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags) error {
	if exists := dependencies.Config.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + dependencies.Config.Runtime.ROOT)
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
		return errors.New("running fmt in root is not supported")
	}

	modules, err := dependencies.Config.GetModules(targets, exclude)
	if err != nil {
		return err
	}

	script := "if [ -n \"$(gofmt -l .)\" ]; then exit 1; fi"

	for _, module := range modules {
		path := filepath.Join(dependencies.Config.Runtime.ROOT, module.RelativePath)
		if err := dependencies.Effects.Executor.Bash(path, script); err != nil {
			return errors.New("error: fmt-ci failed in module " + module.Name)
		}
	}

	return nil
}
