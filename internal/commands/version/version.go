package version

import (
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
)

// Version is injected at build time
var Version = "dev"

func version(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags) error {
	dependencies.Effects.Logger.DefaultLn(Version)
	return nil
}
