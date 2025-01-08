package exec

import (
	"errors"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"path/filepath"
	"strconv"
	"strings"
)

func exec(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags, scriptName string) error {
	if exists := dependencies.Config.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + dependencies.Config.Runtime.ROOT)
	}

	if globalFlags.Verbose {
		dependencies.Effects.Logger.VerboseLn("verbose mode enabled")
	}

	if scriptName == "" {
		return errors.New("no script name provided, usage: gorepo run [script_name]")
	} else {
		if globalFlags.Verbose {
			dependencies.Effects.Logger.VerboseLn("running script '" + scriptName + "'")
		}
	}

	if globalFlags.Verbose {
		dependencies.Effects.Logger.VerboseLn("value for flag allowMissing: " + strconv.FormatBool(cmdFlags.AllowMissing))
	}

	targets := strings.Split(cmdFlags.Target, ",")
	if globalFlags.Verbose {
		dependencies.Effects.Logger.VerboseLn("value for flag target:       " + cmdFlags.Target)
	}

	exclude := strings.Split(cmdFlags.Exclude, ",")
	if globalFlags.Verbose {
		dependencies.Effects.Logger.VerboseLn("value for flag exclude:      " + cmdFlags.Exclude)
	}

	if targets[0] == "root" {
		rootConfig, err := dependencies.Config.GetRootConfig()
		if err != nil {
			return err
		}
		path := dependencies.Config.Runtime.ROOT
		script := rootConfig.Scripts[scriptName]
		if script == "" {
			dependencies.Effects.Logger.InfoLn("script is empty, skipping")
			return errors.New("There is no script named " + scriptName + " at root")
		}
		dependencies.Effects.Logger.InfoLn("running script " + scriptName + " at root ")
		if err := dependencies.Effects.Executor.Bash(path, script); err != nil {
			return err
		}
		return nil
	}

	modules, err := dependencies.Config.GetModules(targets, exclude)
	if err != nil {
		return err
	}

	if len(modules) == 0 {
		return errors.New("no modules found")
	}

	// check all modules have the script
	if globalFlags.Verbose && !cmdFlags.AllowMissing {
		dependencies.Effects.Logger.VerboseLn("checking if all modules have the script")
	}
	var modulesWithoutScript []string
	for _, module := range modules {
		if _, ok := module.Scripts[scriptName]; !ok || module.Scripts[scriptName] == "" {
			modulesWithoutScript = append(modulesWithoutScript, module.Name)
		}
	}
	if len(modulesWithoutScript) == len(modules) {
		return errors.New("not running script, because it is missing in all modules")
	} else if len(modulesWithoutScript) > 0 && !cmdFlags.AllowMissing {
		return errors.New("not running script, because it is missing in following modules '" + scriptName + "' :" + strings.Join(modulesWithoutScript, ", "))
	} else if len(modulesWithoutScript) > 0 && cmdFlags.AllowMissing {
		if globalFlags.Verbose {
			dependencies.Effects.Logger.VerboseLn("script is missing in following modules (but flag allowMissing was passed) '" + scriptName + "' :" + strings.Join(modulesWithoutScript, ", "))
		}
	} else {
		if globalFlags.Verbose {
			dependencies.Effects.Logger.VerboseLn("all modules have the script")
		}
	}

	// execute them
	for _, module := range modules {
		path := filepath.Join(dependencies.Config.Runtime.ROOT, module.RelativePath)
		script := module.Scripts[scriptName]
		if script == "" {
			dependencies.Effects.Logger.InfoLn("script is empty, skipping")
			continue
		}
		dependencies.Effects.Logger.InfoLn("running script " + scriptName + " in module " + module.Name)
		if err := dependencies.Effects.Executor.Bash(path, script); err != nil {
			return err
		}
	}

	return nil
}
