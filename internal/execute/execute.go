package execute

import (
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/commands"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"gorepo-cli/pkg"
	"gorepo-cli/pkg/logger"
	"os"
)

func Execute(l logger.Methods) (err error) {
	effects := pkg.NewEffects(l)
	cfg, err := config.NewConfig(effects)
	if err != nil {
		return err
	}
	dependencies := config.NewDependencies(effects, cfg)

	app := &cli.App{
		Name:     "gorepo",
		Usage:    "A cli tool to manage Go monorepos",
		Commands: commands.RegisterCommands(dependencies),
		Flags:    flags.GlobalGroup,
	}

	if exists := cfg.RootConfigExists(); exists == true {
		modules, err := cfg.GetModules([]string{"all"}, nil)
		if err != nil {
			return err
		}

		for _, module := range modules {
			module := module
			app.Commands = append(app.Commands, &cli.Command{
				Name:        module.Name,
				Hidden:      true,
				Subcommands: commands.RegisterModuleCommands(module.Name, dependencies),
			})
		}
	}

	return app.Run(os.Args)
}
