package initcmd

import (
	"errors"
	"fmt"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"path/filepath"
)

func initCmd(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags, name string) error {
	if exists := dependencies.Config.RootConfigExists(); exists {
		return errors.New("monorepo already exists at " + dependencies.Config.Runtime.ROOT)
	}

	verbose := globalFlags.Verbose

	rootConfig := config.RootConfig{
		Name:    name,
		Version: "0.1.0",
		Vendor:  true,
	}

	// ask name
	if rootConfig.Name == "" {
		defaultName := filepath.Base(dependencies.Config.Runtime.ROOT)
		nameResponse, err := dependencies.Effects.Terminal.AskString("What is the monorepo name?", "", defaultName, dependencies.Effects.Logger)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		rootConfig.Name = nameResponse
	}

	// ask if should vendor
	if vendorResponse, err := dependencies.Effects.Terminal.AskBool(
		"Do you want to vendor dependencies?", "y/n", "y", dependencies.Effects.Logger); err == nil {
		rootConfig.Vendor = vendorResponse
	} else {
		return fmt.Errorf("failed to read input: %w", err)
	}

	if exists := dependencies.Config.GoWorkspaceExists(); !exists {
		if verbose {
			dependencies.Effects.Logger.VerboseLn("go workspace does not exist yet, running 'go work init'")
		}
		err := dependencies.Effects.Executor.Go(dependencies.Config.Runtime.ROOT, "work", "init")
		if err != nil {
			return err
		}
	} else {
		if verbose {
			dependencies.Effects.Logger.VerboseLn("go workspace already exists, no need to create one")
		}
		// todo: handle vendoring
	}

	if err := dependencies.Config.WriteRootConfig(rootConfig); err != nil {
		return err
	} else {
		if verbose {
			dependencies.Effects.Logger.VerboseLn("created monorepo configuration 'work.toml' at root")
		}
	}

	// todo: check existence of modules folder (go.mod) to sanitize everything (create module.toml and make sure they are in the workspace)

	dependencies.Effects.Logger.SuccessLn("monorepo successfully initialized at " + dependencies.Config.Runtime.ROOT)

	return nil
}
