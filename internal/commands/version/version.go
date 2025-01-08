package version

import (
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
)

// _version is injected at build time
var _version = "dev"

func version(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags) error {
	dependencies.Effects.Logger.DefaultLn(_version)
	return nil
}
