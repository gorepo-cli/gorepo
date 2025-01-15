package reusable

import (
	"github.com/pelletier/go-toml/v2"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"gorepo-cli/pkg"
	"testing"
)

func TestExecuteGenericGoCommand(t *testing.T) {
	t.Run("should run command if no task is provided", func(t *testing.T) {
		t.Run("at root", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootNoTask),
					"/some/path/root/mod1/module.toml": []byte(goModuleNoTask),
					"/some/path/root/mod2/module.toml": []byte(goModuleNoTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := ExecuteGenericGoCommand(dependencies, &flags.CommandFlags{
				Target: "root",
			}, &flags.GlobalFlags{}, "test_command", "test_script")
			if err != nil {
				t.Fatal(err)
			}
			commands := effects.Executor.Output()
			expected := "go " + "test_command"
			if len(commands) != 2 {
				t.Fatalf("Expected %d, got %d", 2, len(commands))
			}
			if commands[0].Dir != "/some/path/root/mod1" {
				t.Fatalf("Expected %s, got %s", "/some/path/root/mod1", commands[0].Dir)
			}
			if commands[0].Command != expected {
				t.Fatalf("Expected %s, got %s", expected, commands[0].Command)
			}
			if commands[1].Dir != "/some/path/root/mod2" {
				t.Fatalf("Expected %s, got %s", "/some/path/root/mod2", commands[1].Dir)
			}
			if commands[1].Command != expected {
				t.Fatalf("Expected %s, got %s", expected, commands[1].Command)
			}
		})
		t.Run("in go modules", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootNoTask),
					"/some/path/root/mod1/module.toml": []byte(goModuleNoTask),
					"/some/path/root/mod2/module.toml": []byte(goModuleNoTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := ExecuteGenericGoCommand(dependencies, &flags.CommandFlags{
				Target: "all",
			}, &flags.GlobalFlags{}, "test_command", "test_script")
			if err != nil {
				t.Fatal(err)
			}
			expected := "go " + "test_command"
			commands := effects.Executor.Output()
			if len(commands) != 2 {
				t.Fatalf("Expected %d, got %d", 2, len(commands))
			}
			if commands[0].Dir != "/some/path/root/mod1" {
				t.Fatalf("Expected %s, got %s", "/some/path/root/mod1", commands[0].Dir)
			}
			if commands[1].Dir != "/some/path/root/mod2" {
				t.Fatalf("Expected %s, got %s", "/some/path/root/mod2", commands[1].Dir)
			}
			if commands[0].Command != expected {
				t.Fatalf("Expected %s, got %s", expected, commands[0].Command)
			}
			if commands[1].Command != expected {
				t.Fatalf("Expected %s, got %s", expected, commands[1].Command)
			}
		})
		t.Run("in non-go modules", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootNoTask),
					"/some/path/root/mod1/module.toml": []byte(nonGoModuleNoTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := ExecuteGenericGoCommand(dependencies, &flags.CommandFlags{
				Target: "all",
			}, &flags.GlobalFlags{}, "test_command", "test_script")
			if err != nil {
				t.Fatal(err)
			}
			commands := effects.Executor.Output()
			if len(commands) != 0 {
				t.Fatalf("Expected %d, got %d", 0, len(commands))
			}
		})
	})
	t.Run("should run task if a task if provided", func(t *testing.T) {
		t.Run("should run task with command name at root", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootWithTask),
					"/some/path/root/mod1/module.toml": []byte(goModuleWithTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := ExecuteGenericGoCommand(dependencies, &flags.CommandFlags{
				Target: "root",
			}, &flags.GlobalFlags{}, "test_command", "test_script")
			if err != nil {
				t.Fatal(err)
			}
			expected := "bash " + "some_script_from_command"
			commands := effects.Executor.Output()
			if len(commands) != 1 {
				t.Fatalf("Expected %d, got %d", 1, len(commands))
			}
			if commands[0].Dir != "/some/path/root" {
				t.Fatalf("Expected %s, got %s", "/some/path/root", commands[0].Dir)
			}
			if commands[0].Command != expected {
				t.Fatalf("Expected %s, got %s", expected, commands[0].Command)
			}
		})
		t.Run("should run task with command name in go modules", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootWithTask),
					"/some/path/root/mod1/module.toml": []byte(goModuleWithTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := ExecuteGenericGoCommand(dependencies, &flags.CommandFlags{
				Target: "mod1",
			}, &flags.GlobalFlags{}, "test_command", "test_script")
			if err != nil {
				t.Fatal(err)
			}
			expected := "bash " + "some_script_from_command"
			commands := effects.Executor.Output()
			if len(commands) != 1 {
				t.Fatalf("Expected %d, got %d", 1, len(commands))
			}
			if commands[0].Dir != "/some/path/root" {
				t.Fatalf("Expected %s, got %s", "/some/path/root/mod1", commands[0].Dir)
			}
			if commands[0].Command != expected {
				t.Fatalf("Expected %s, got %s", expected, commands[0].Command)
			}
		})
		t.Run("should run task with command name in non-go modules with a task", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootWithTask),
					"/some/path/root/mod1/module.toml": []byte(nonGoModuleWithTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := ExecuteGenericGoCommand(dependencies, &flags.CommandFlags{
				Target: "mod1",
			}, &flags.GlobalFlags{}, "test_command", "test_script")
			if err != nil {
				t.Fatal(err)
			}
			expected := "bash " + "some_script_from_command"
			commands := effects.Executor.Output()
			if len(commands) != 1 {
				t.Fatalf("Expected %d, got %d", 1, len(commands))
			}
			if commands[0].Dir != "/some/path/root" {
				t.Fatalf("Expected %s, got %s", "/some/path/root/mod1", commands[0].Dir)
			}
			if commands[0].Command != expected {
				t.Fatalf("Expected %s, got %s", expected, commands[0].Command)
			}
		})
		t.Run("should not run anything in non-go modules with no task", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootWithTask),
					"/some/path/root/mod1/module.toml": []byte(nonGoModuleNoTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := ExecuteGenericGoCommand(dependencies, &flags.CommandFlags{
				Target: "mod1",
			}, &flags.GlobalFlags{}, "test_command", "test_script")
			if err != nil {
				t.Fatal(err)
			}
			commands := effects.Executor.Output()
			if len(commands) != 0 {
				t.Fatalf("Expected %d, got %d", 1, len(commands))
			}
		})
	})
	t.Run("should run queues of tasks if provided", func(t *testing.T) {
		t.Run("at root", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootWithTaskQueue),
					"/some/path/root/mod1/module.toml": []byte(goModuleNoTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := ExecuteGenericGoCommand(dependencies, &flags.CommandFlags{
				Target: "root",
			}, &flags.GlobalFlags{}, "test_command", "test_script")
			if err != nil {
				t.Fatal(err)
			}
			expected := "bash " + "some_script_from_command"
			commands := effects.Executor.Output()
			if len(commands) != 2 {
				t.Fatalf("Expected %d, got %d", 2, len(commands))
			}
			if commands[0].Dir != "/some/path/root" {
				t.Fatalf("Expected %s, got %s", "/some/path/root", commands[0].Dir)
			}
			if commands[0].Command != expected+"_1" {
				t.Fatalf("Expected %s, got %s", expected+"_1", commands[0].Command)
			}
			if commands[1].Dir != "/some/path/root" {
				t.Fatalf("Expected %s, got %s", "/some/path/root", commands[1].Dir)
			}
			if commands[1].Command != expected+"_2" {
				t.Fatalf("Expected %s, got %s", expected+"_2", commands[1].Command)
			}
		})
		t.Run("in module", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootNoTask),
					"/some/path/root/mod1/module.toml": []byte(nonGoModuleWithTaskQueue),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := ExecuteGenericGoCommand(dependencies, &flags.CommandFlags{
				Target: "mod1",
			}, &flags.GlobalFlags{}, "test_command", "test_script")
			if err != nil {
				t.Fatal(err)
			}
			expected := "bash " + "some_script_from_command"
			commands := effects.Executor.Output()
			if len(commands) != 2 {
				t.Fatalf("Expected %d, got %d", 2, len(commands))
			}
			if commands[0].Dir != "/some/path/root" {
				t.Fatalf("Expected %s, got %s", "/some/path/root/mod1", commands[0].Dir)
			}
			if commands[0].Command != expected+"_1" {
				t.Fatalf("Expected %s, got %s", expected+"_1", commands[0].Command)
			}
			if commands[1].Dir != "/some/path/root" {
				t.Fatalf("Expected %s, got %s", "/some/path/root/mod2", commands[1].Dir)
			}
			if commands[1].Command != expected+"_2" {
				t.Fatalf("Expected %s, got %s", expected+"_2", commands[1].Command)
			}
		})
	})
	t.Run("should target modules when using --target and --exclude", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml":        []byte(rootNoTask),
				"/some/path/root/mod1/module.toml": []byte(goModuleNoTask),
				"/some/path/root/mod2/module.toml": []byte(goModuleNoTask),
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := ExecuteGenericGoCommand(dependencies, &flags.CommandFlags{
			Target:  "all",
			Exclude: "mod1",
		}, &flags.GlobalFlags{}, "test_command", "test_script")
		if err != nil {
			t.Fatal(err)
		}
		expected := "go " + "test_command"
		commands := effects.Executor.Output()
		if len(commands) != 1 {
			t.Fatalf("Expected %d, got %d", 1, len(commands))
		}
		if commands[0].Dir != "/some/path/root/mod2" {
			t.Fatalf("Expected %s, got %s", "/some/path/root/mod2", commands[0].Dir)
		}
		if commands[0].Command != expected {
			t.Fatalf("Expected %s, got %s", expected, commands[0].Command)
		}
	})
}

var goModuleNoTask, _ = toml.Marshal(config.ModuleConfig{
	Template: "@default",
	Type:     "any",
	Language: "go",
})

var goModuleWithTask, _ = toml.Marshal(config.ModuleConfig{
	Template: "@default",
	Type:     "any",
	Language: "go",
	Tasks: map[string]config.TaskQueue{
		"test_command": config.TaskQueue{"some_script_from_command"},
	},
})

var nonGoModuleNoTask, _ = toml.Marshal(config.ModuleConfig{
	Template: "@default",
	Type:     "any",
	Language: "other",
})

var nonGoModuleWithTask, _ = toml.Marshal(config.ModuleConfig{
	Template: "@default",
	Type:     "any",
	Language: "other",
	Tasks: map[string]config.TaskQueue{
		"test_command": config.TaskQueue{"some_script_from_command"},
	},
})

var nonGoModuleWithTaskQueue, _ = toml.Marshal(config.ModuleConfig{
	Template: "@default",
	Type:     "any",
	Language: "other",
	Tasks: map[string]config.TaskQueue{
		"test_command": config.TaskQueue{"some_script_from_command_1", "some_script_from_command_2"},
	},
})

var rootNoTask, _ = toml.Marshal(config.RootConfig{
	Name: "test",
})

var rootWithTask, _ = toml.Marshal(config.RootConfig{
	Name: "test",
	Tasks: map[string]config.TaskQueue{
		"test_command": config.TaskQueue{"some_script_from_command"},
	},
})

var rootWithTaskQueue, _ = toml.Marshal(config.RootConfig{
	Name: "test",
	Tasks: map[string]config.TaskQueue{
		"test_command": config.TaskQueue{"some_script_from_command_1", "some_script_from_command_2"},
	},
})
