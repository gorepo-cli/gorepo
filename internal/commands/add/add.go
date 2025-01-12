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
	if relativePathAndNameInput == "" {
		return errors.New("error: no name provided")
	}
	name := filepath.Base(relativePathAndNameInput)
	if globalFlags.Verbose {
		dependencies.Effects.Logger.VerboseLn("the relative path is " + relativePathAndNameInput + " and the name is " + name)
	}
	modules, err := dependencies.Config.GetModules([]string{"all"}, []string{})
	if err != nil {
		return err
	}
	for _, module := range modules {
		if module.Name == name {
			return errors.New("the module with name " + name + " already exists at the path " + module.RelativePath)
		}
	}
	moduleType, err := dependencies.Effects.Terminal.SingleSelect(
		"what type of module do you want to create?",
		[][]string{
			{"executable", "meant to be built and executed"},
			{"library   ", "meant to be built, not executed"},
			{"script    ", "meant to be executed, not built"},
			{"static    ", "meant to be imported directly"},
		},
		dependencies.Effects.Logger,
	)
	if err != nil {
		return err
	}
	language, err := dependencies.Effects.Terminal.SingleSelect(
		"what language is it using?",
		[][]string{
			{"go        ", ""},
			{"javascript", ""},
			{"other     ", ""},
		},
		dependencies.Effects.Logger,
	)
	if err != nil {
		return err
	}
	newModule := config.ModuleConfig{
		Name:         name,
		RelativePath: relativePathAndNameInput,
		Template:     "@default",
		Type:         moduleType,
		Language:     language,
		Main:         "",
		Priority:     0,
		Scripts:      map[string]config.ScriptQueue{},
	}
	absolutePath := filepath.Join(dependencies.Config.Runtime.ROOT, relativePathAndNameInput)
	if err := dependencies.Config.WriteModuleConfig(newModule, absolutePath); err != nil {
		return err
	}
	if err := dependencies.Effects.Executor.Go(absolutePath, "mod", "init", name); err != nil {
		return err
	}
	if language == "go" {
		if err := dependencies.Effects.Executor.Go(dependencies.Config.Runtime.ROOT, "work", "use", relativePathAndNameInput); err != nil {
			return err
		}
	}
	return nil
}
