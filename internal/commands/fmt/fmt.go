package fmt

import (
	"errors"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"path/filepath"
	"strings"
)

func fmt(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags) error {
	if err := dependencies.Config.BreakIfRootConfigDoesNotExist(); err != nil {
		return err
	}

	if globalFlags.Verbose {
		dependencies.Effects.Logger.VerboseLn("the value of --target is '" + cmdFlags.Target + "'")
		dependencies.Effects.Logger.VerboseLn("the value of --exclude is '" + cmdFlags.Exclude + "'")
	}

	targets := strings.Split(cmdFlags.Target, ",")
	exclude := strings.Split(cmdFlags.Exclude, ",")

	// executing the script at the root means executing it in all modules
	if targets[0] == "root" {
		targets = []string{"all"}
	}

	modules, err := dependencies.Config.GetModules(targets, exclude)
	if err != nil {
		return err
	}

	script := "if [ -n \"$(gofmt -l .)\" ]; then exit 1; fi"

	for _, module := range modules {
		path := filepath.Join(dependencies.Config.Runtime.ROOT, module.RelativePath)
		if module.Language == "go" {
			if cmdFlags.Ci == true {
				if err := dependencies.Effects.Executor.Bash(path, script); err != nil {
					return errors.New("error: fmt failed in module " + module.Name)
				}
			} else {
				if err := dependencies.Effects.Executor.Go(path, "fmt"); err != nil {
					return errors.New("error: fmt failed in module " + module.Name)
				}
			}
		} else {
			if cmdFlags.Ci == true {
				hasFmtCi := false
				for name, _ := range module.Scripts {
					if name == "fmt-ci" {
						hasFmtCi = true
						break
					}
				}
				if hasFmtCi {
					dependencies.Effects.Logger.InfoLn("module " + module.Name + " in not a go module but it implements fmt, executing")
					for i, _ := range module.Scripts["fmt-ci"] {
						if err := dependencies.Effects.Executor.Bash(path, module.Scripts["fmt-ci"][i]); err != nil {
							return errors.New("error: fmt-ci failed in module " + module.Name)
						}
					}
				} else {
					dependencies.Effects.Logger.WarningLn("module " + module.Name + " in not a go module and it does not implement fmt, skipping")
				}
			} else {
				hasFmt := false
				for name, _ := range module.Scripts {
					if name == "fmt" {
						hasFmt = true
						break
					}
				}
				if hasFmt {
					dependencies.Effects.Logger.InfoLn("module " + module.Name + " in not a go module but it implements fmt, executing")
					for i, _ := range module.Scripts["fmt"] {
						if err := dependencies.Effects.Executor.Bash(path, module.Scripts["fmt"][i]); err != nil {
							return errors.New("error: fmt failed in module " + module.Name)
						}
					}
				} else {
					dependencies.Effects.Logger.WarningLn("module " + module.Name + " in not a go module and it does not implement fmt, skipping")
				}
			}
		}
	}

	return nil
}
