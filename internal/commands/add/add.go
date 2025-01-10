package add

import (
	"errors"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"path/filepath"
)

func add(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags, relativePathAndNameInput string) error {
	if exists := dependencies.Config.RootConfigExists(); !exists {
		return errors.New("monorepo not found at " + dependencies.Config.Runtime.ROOT)
	}
	dependencies.Effects.Logger.VerboseLn("relativePathAndNameInput: " + relativePathAndNameInput)
	if relativePathAndNameInput == "" {
		return errors.New("error: no name provided")
	}
	name := filepath.Base(relativePathAndNameInput)
	dependencies.Effects.Logger.VerboseLn("name: " + name)
	if modules, err := dependencies.Config.GetModules([]string{"all"}, []string{}); err != nil {
		return err
	} else {
		for _, module := range modules {
			if module.Name == name {
				return errors.New("module with name " + name + " already exists at " + module.RelativePath)
			}
		}
	}
	newModule := config.ModuleConfig{
		Name:         name,
		RelativePath: relativePathAndNameInput,
		Template:     "@default",
		Type:         "executable",
		Main:         "",
		Priority:     0,
		Scripts:      map[string]config.Pipeline{},
	}
	absolutePath := filepath.Join(dependencies.Config.Runtime.ROOT, relativePathAndNameInput)
	if err := dependencies.Config.WriteModuleConfig(newModule, absolutePath); err != nil {
		return err
	}
	if err := dependencies.Effects.Executor.Go(absolutePath, "mod", "init", name); err != nil {
		return err
	}
	if rootConfig, err := dependencies.Config.GetRootConfig(); err != nil {
		return err
	} else if rootConfig.Strategy == "workspace" {
		if err := dependencies.Effects.Executor.Go(dependencies.Config.Runtime.ROOT, "work", "use", relativePathAndNameInput); err != nil {
			return err
		}
	}
	return nil
}
