package fmt

import (
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"gorepo-cli/internal/reusable"
)

func fmt(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags) error {
	if err := dependencies.Config.BreakIfRootConfigDoesNotExist(); err != nil {
		return err
	}

	if cmdFlags.Ci == true {
		script := "if [ -n \"$(gofmt -l .)\" ]; then exit 1; fi"
		return reusable.ExecuteGenericGoCommand(dependencies, cmdFlags, globalFlags, "fmt-ci", script)
	}

	return reusable.ExecuteGenericGoCommand(dependencies, cmdFlags, globalFlags, "fmt", "fmt")
}
