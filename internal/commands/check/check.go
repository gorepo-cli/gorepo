package check

import (
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"strconv"
)

func check(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags) error {
	dependencies.Effects.Logger.InfoLn("===================")
	dependencies.Effects.Logger.InfoLn("RUNTIME_CONFIG")
	dependencies.Effects.Logger.InfoLn("===================")
	dependencies.Effects.Logger.DefaultLn("WD_(COMMAND_RAN_FROM)........" + dependencies.Config.Runtime.WD)
	dependencies.Effects.Logger.DefaultLn("ROOT (OF THE MONOREPO)......." + dependencies.Config.Runtime.ROOT)
	dependencies.Effects.Logger.DefaultLn("MONOREPO EXISTS (AT ROOT)...." +
		strconv.FormatBool(dependencies.Config.RootConfigExists()))

	dependencies.Effects.Logger.InfoLn("===================")
	dependencies.Effects.Logger.InfoLn("STATIC_CONFIG")
	dependencies.Effects.Logger.InfoLn("===================")
	dependencies.Effects.Logger.DefaultLn("MAX RECURSION................" + strconv.Itoa(dependencies.Config.Static.MaxRecursion))
	dependencies.Effects.Logger.DefaultLn("ROOT FILE NAME..............." + dependencies.Config.Static.RootFileName)
	dependencies.Effects.Logger.DefaultLn("MODULE FILE NAME............." + dependencies.Config.Static.ModuleFileName)

	if dependencies.Config.RootConfigExists() {
		dependencies.Effects.Logger.InfoLn("===================")
		dependencies.Effects.Logger.InfoLn("ROOT_CONFIG")
		dependencies.Effects.Logger.InfoLn("===================")

		cfg, err := dependencies.Config.GetRootConfig()
		if err != nil {
			return err
		}

		dependencies.Effects.Logger.DefaultLn("NAME.........." + cfg.Name)
		dependencies.Effects.Logger.DefaultLn("VERSION......." + cfg.Version)
		dependencies.Effects.Logger.DefaultLn("STRATEGY......" + cfg.Strategy)
		dependencies.Effects.Logger.DefaultLn("VENDOR........" + strconv.FormatBool(cfg.Vendor))

		modules, err := dependencies.Config.GetModules([]string{"all"}, []string{})
		if err != nil {
			return err
		}

		dependencies.Effects.Logger.DefaultLn("N_MODULES....." + strconv.Itoa(len(modules)))

		if len(modules) > 0 {
			dependencies.Effects.Logger.InfoLn("===================")
			dependencies.Effects.Logger.InfoLn("MODULES_CONFIG")
			dependencies.Effects.Logger.InfoLn("===================")
		}

		for _, module := range modules {
			dependencies.Effects.Logger.InfoLn("MODULE " + module.Name)
			dependencies.Effects.Logger.DefaultLn("MODULE_NAME........ " + module.Name)
			dependencies.Effects.Logger.DefaultLn("MODULE_PATH........ " + module.RelativePath)
			if len(module.Scripts) > 0 {
				dependencies.Effects.Logger.DefaultLn("COMMANDS........")
				for k, v := range module.Scripts {
					dependencies.Effects.Logger.DefaultLn("  " + k + " -> " + v)
				}
			}
		}
	}
	return nil
}
