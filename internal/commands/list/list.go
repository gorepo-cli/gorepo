package list

import (
	"errors"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
)

func list(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags) error {
	if exists := dependencies.Config.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + dependencies.Config.Runtime.ROOT)
	}
	modules, err := dependencies.Config.GetModules([]string{"all"}, []string{})
	if err != nil {
		return err
	}
	if len(modules) == 0 {
		dependencies.Effects.Logger.InfoLn("no modules found")
	} else {
		for _, module := range modules {
			dependencies.Effects.Logger.DefaultLn(module.Name)
		}
	}
	return nil
}
