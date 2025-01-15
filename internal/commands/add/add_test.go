package add

import (
	"github.com/pelletier/go-toml/v2"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"gorepo-cli/pkg"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Run("should fail if there is no work.toml", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD:    "/some/path/root",
			Files: map[string][]byte{},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := add(dependencies, &flags.CommandFlags{
			Target: "all",
		}, &flags.GlobalFlags{}, "mod1")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
	t.Run("should fail if the name is already taken (at root)", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml":        rootNoTask,
				"/some/path/root/mod1/module.toml": moduleNoTask,
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := add(dependencies, &flags.CommandFlags{}, &flags.GlobalFlags{}, "mod1")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
	t.Run("should fail if the name is already taken (at path)", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml":        rootNoTask,
				"/some/path/root/mod1/module.toml": moduleNoTask,
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := add(dependencies, &flags.CommandFlags{}, &flags.GlobalFlags{}, "path/mod1")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
	t.Run("should generate a new module", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml":        rootNoTask,
				"/some/path/root/mod1/module.toml": moduleNoTask,
			},
			QaSingleSelect: map[string]string{
				"what type of module do you want to create?": "some_type",
				"what language is it using?":                 "some_language",
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := add(dependencies, &flags.CommandFlags{}, &flags.GlobalFlags{}, "mod2")
		if err != nil {
			t.Fatal(err)
		}
		files := effects.Filesystem.Output()
		if len(files) != 3 {
			t.Fatalf("Expected %d, got %d", 3, len(files))
		}
		newModuleRaw := files["/some/path/root/mod2/module.toml"]
		if newModuleRaw == nil {
			t.Fatalf("Expected %s, got nil", "a file")
		}
		var newModule config.ModuleConfig
		err = toml.Unmarshal(newModuleRaw, &newModule)
		if newModule.Template != "@default" {
			t.Fatalf("Expected %s, got %s", "@default", newModule.Template)
		}
		if newModule.Language != "some_language" {
			t.Fatalf("Expected %s, got %s", "some_language", newModule.Language)
		}
		if newModule.Type != "some_type" {
			t.Fatalf("Expected %s, got %s", "some_type", newModule.Type)
		}
	})
	t.Run("should generate a new module from root, even if called from another location", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root/mod1",
			Files: map[string][]byte{
				"/some/path/root/work.toml":        rootNoTask,
				"/some/path/root/mod1/module.toml": moduleNoTask,
			},
			QaSingleSelect: map[string]string{
				"what type of module do you want to create?": "some_type",
				"what language is it using?":                 "some_language",
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := add(dependencies, &flags.CommandFlags{}, &flags.GlobalFlags{}, "mod2")
		if err != nil {
			t.Fatal(err)
		}
		files := effects.Filesystem.Output()
		if len(files) != 3 {
			t.Fatalf("Expected %d, got %d", 3, len(files))
		}
		newModuleRaw := files["/some/path/root/mod2/module.toml"]
		if newModuleRaw == nil {
			t.Fatalf("Expected %s, got nil", "a file")
		}
		var newModule config.ModuleConfig
		err = toml.Unmarshal(newModuleRaw, &newModule)
		if newModule.Template != "@default" {
			t.Fatalf("Expected %s, got %s", "@default", newModule.Template)
		}
		if newModule.Language != "some_language" {
			t.Fatalf("Expected %s, got %s", "some_language", newModule.Language)
		}
		if newModule.Type != "some_type" {
			t.Fatalf("Expected %s, got %s", "some_type", newModule.Type)
		}
		commands := effects.Executor.Output()
		expected := "go mod init mod2"
		if len(commands) != 1 {
			t.Fatalf("Expected %d, got %d", 1, len(commands))
		}
		if commands[0].Dir != "/some/path/root/mod2" {
			t.Fatalf("Expected %s, got %s", "/some/path/root", commands[0].Dir)
		}
		if commands[0].Command != expected {
			t.Fatalf("Expected %s, got %s", expected, commands[0].Command)
		}
	})
	t.Run("should add go module to go.work if the language is go", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml":        rootNoTask,
				"/some/path/root/mod1/module.toml": moduleNoTask,
			},
			QaSingleSelect: map[string]string{
				"what type of module do you want to create?": "some_type",
				"what language is it using?":                 "go",
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := add(dependencies, &flags.CommandFlags{}, &flags.GlobalFlags{}, "mod2")
		if err != nil {
			t.Fatal(err)
		}
		commands := effects.Executor.Output()
		expected := "go work use mod2"
		if len(commands) != 2 {
			t.Fatalf("Expected %d, got %d", 2, len(commands))
		}
		if commands[1].Dir != "/some/path/root" {
			t.Fatalf("Expected %s, got %s", "/some/path/root", commands[0].Dir)
		}
		if commands[1].Command != expected {
			t.Fatalf("Expected %s, got %s", expected, commands[0].Command)
		}
	})
	t.Run("should not add go module to go.work if the language is not go", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml":        rootNoTask,
				"/some/path/root/mod1/module.toml": moduleNoTask,
			},
			QaSingleSelect: map[string]string{
				"what type of module do you want to create?": "some_type",
				"what language is it using?":                 "other",
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := add(dependencies, &flags.CommandFlags{}, &flags.GlobalFlags{}, "mod2")
		if err != nil {
			t.Fatal(err)
		}
		commands := effects.Executor.Output()
		if len(commands) != 1 {
			t.Fatalf("Expected %d, got %d", 1, len(commands))
		}
	})
}

var rootNoTask, _ = toml.Marshal(config.RootConfig{
	Name: "test",
})

var moduleNoTask, _ = toml.Marshal(config.ModuleConfig{
	Template: "@default",
	Type:     "any",
	Language: "go",
})
