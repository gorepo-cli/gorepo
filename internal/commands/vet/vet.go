package vet

import (
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"gorepo-cli/internal/reusable"
)

func vet(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags) error {
	if err := dependencies.Config.BreakIfRootConfigDoesNotExist(); err != nil {
		return err
	}

	if cmdFlags.Ci == true {
		script := "go vet . || exit 1"
		return reusable.ExecuteGenericGoCommand(dependencies, cmdFlags, globalFlags, "vet-ci", script)
	}

	return reusable.ExecuteGenericGoCommand(dependencies, cmdFlags, globalFlags, "vet", "vet")
}
