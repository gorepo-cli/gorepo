package exec

import (
	"errors"
	"fmt"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"path/filepath"
	"strconv"
	"strings"
)

func exec(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags, scriptName string) error {
	if exists := dependencies.Config.RootConfigExists(); !exists {
		return errors.New(fmt.Sprintf("no monorepo found at path %s" + dependencies.Config.Runtime.ROOT))
	}

	if globalFlags.Verbose {
		dependencies.Effects.Logger.VerboseLn("the script to run is '" + scriptName + "'")
		dependencies.Effects.Logger.VerboseLn("the value of --target is '" + cmdFlags.Target + "'")
		dependencies.Effects.Logger.VerboseLn("the value of --exclude is '" + cmdFlags.Exclude + "'")
		dependencies.Effects.Logger.VerboseLn("the value of --allow-missing is '" + strconv.FormatBool(cmdFlags.AllowMissing) + "'")
	}

	if scriptName == "" {
		return errors.New("no script provided, type gorepo exec help to see usage")
	}

	targets := strings.Split(cmdFlags.Target, ",")
	exclude := strings.Split(cmdFlags.Exclude, ",")

	if targets[0] == "root" {
		rootConfig, err := dependencies.Config.GetRootConfig()
		if err != nil {
			return err
		}

		path := dependencies.Config.Runtime.ROOT
		script := rootConfig.Scripts[scriptName]
		if len(script) == 0 {
			return errors.New(fmt.Sprintf("there is no script named '%s' at root", scriptName))
		}

		for i, unitScript := range script {
			if len(script) == 1 {
				dependencies.Effects.Logger.InfoLn("running script '" + scriptName + "' at root")
			} else {
				dependencies.Effects.Logger.InfoLn("running script '" + scriptName + "', step " + strconv.Itoa(i+1) + " at root")
			}
			if err := dependencies.Effects.Executor.Bash(path, unitScript); err != nil {
				if len(script) > 1 {
					dependencies.Effects.Logger.WarningLn(fmt.Sprintf("/!\\ script '%s' failed at root, at step %d, be aware earlier steps may have run", scriptName, i+1))
				}
				return err
			}
		}

		return nil
	}

	modules, err := dependencies.Config.GetModules(targets, exclude)
	if err != nil {
		return err
	}

	if len(modules) == 0 {
		return errors.New("no module matches the target, use --verbose to debug")
	} else if globalFlags.Verbose {
		dependencies.Effects.Logger.VerboseLn("=== modules that matches the criteria")
		for _, module := range modules {
			dependencies.Effects.Logger.VerboseLn("--- " + module.Name)
		}
	}

	var modulesWithoutScript []string

	for _, module := range modules {
		if _, ok := module.Scripts[scriptName]; !ok || module.Scripts[scriptName] == nil || len(module.Scripts[scriptName]) == 0 {
			modulesWithoutScript = append(modulesWithoutScript, module.Name)
		}
	}

	if globalFlags.Verbose {
		dependencies.Effects.Logger.VerboseLn("=== modules that does not have the script")
		for _, moduleName := range modulesWithoutScript {
			dependencies.Effects.Logger.VerboseLn("--- " + moduleName)
		}
	}

	if len(modulesWithoutScript) == len(modules) {
		return errors.New("the script is missing in all modules")
	}

	if len(modulesWithoutScript) > 0 && !cmdFlags.AllowMissing {
		return errors.New(fmt.Sprintf(
			"the script is missing in modules %s. Use --allow-missing or --exclude",
			strings.Join(modulesWithoutScript, ", ")))
	}

	var nSuccess = 0
	var nSkipped = 0

	for _, module := range modules {
		path := filepath.Join(dependencies.Config.Runtime.ROOT, module.RelativePath)
		script := module.Scripts[scriptName]

		if len(script) == 0 {
			dependencies.Effects.Logger.WarningLn(fmt.Sprintf("the blue llama is skipping module '%s' (no script '%s')", module.Name, scriptName))
			nSkipped++
			continue
		}

		for i, unitScript := range script {
			if len(script) == 1 {
				dependencies.Effects.Logger.InfoLn(fmt.Sprintf("the blue llama is running script '%s' in module '%s'", scriptName, module.Name))
			} else {
				dependencies.Effects.Logger.InfoLn(fmt.Sprintf("the blue llama is running script '%s', step %d, in module '%s'", scriptName, i+1, module.Name))
			}
			if err := dependencies.Effects.Executor.Bash(path, unitScript); err != nil {
				if len(script) == 1 {
					if len(modules) > 1 {
						dependencies.Effects.Logger.WarningLn(fmt.Sprintf("/!\\ script failed within module '%s', be aware it may have run for other modules", module.Name))
					}
				} else {
					dependencies.Effects.Logger.WarningLn(fmt.Sprintf("/!\\ script failed within module '%s', step %d, be aware it may have run for other modules", module.Name, i+1))
				}
				return err
			}
		}

		nSuccess++
	}

	dependencies.Effects.Logger.SuccessLn(fmt.Sprintf("the blue llama gracefully finished, targeted %s modules, executed %s, skipped %s", strconv.Itoa(len(modules)), strconv.Itoa(nSuccess), strconv.Itoa(nSkipped)))

	return nil
}
