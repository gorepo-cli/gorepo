package exec

import (
	"github.com/pelletier/go-toml/v2"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"gorepo-cli/pkg"
	"testing"
)

func TestExec(t *testing.T) {
	t.Run("general", func(t *testing.T) {
		t.Run("should fail if there is no work.toml", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD:    "/some/path/root",
				Files: map[string][]byte{},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := exec(dependencies, &flags.CommandFlags{
				Target: "all",
			}, &flags.GlobalFlags{}, "script_name")
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
		})
	})
	t.Run("root target", func(t *testing.T) {
		t.Run("should fail if there is no task at root", func(t *testing.T) {
			t.Run("should fail if there is no work.toml", func(t *testing.T) {
				effects := pkg.NewTestkit(pkg.TestKitArgs{
					WD:    "/some/path/root",
					Files: map[string][]byte{},
				})
				cfg, _ := config.NewConfig(effects.ToEffects())
				dependencies := config.NewDependencies(effects.ToEffects(), cfg)
				err := exec(dependencies, &flags.CommandFlags{
					Target: "all",
				}, &flags.GlobalFlags{}, "script_name")
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
			})
		})
		t.Run("should run the task at root", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootWithTask),
					"/some/path/root/mod1/module.toml": []byte(goModuleNoTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := exec(dependencies, &flags.CommandFlags{
				Target: "root",
			}, &flags.GlobalFlags{}, "script_name")
			if err != nil {
				t.Fatal(err)
			}
			commands := effects.Executor.Output()
			expected := "bash " + "some_bash_script"
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
		t.Run("should run task at root, even from a module", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root/mod1",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootWithTask),
					"/some/path/root/mod1/module.toml": []byte(goModuleNoTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := exec(dependencies, &flags.CommandFlags{
				Target: "root",
			}, &flags.GlobalFlags{}, "script_name")
			if err != nil {
				t.Fatal(err)
			}
			commands := effects.Executor.Output()
			expected := "bash " + "some_bash_script"
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
	})
	t.Run("modules target", func(t *testing.T) {
		t.Run("should fail if the task is missing in one module", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root/mod1",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootNoTask),
					"/some/path/root/mod1/module.toml": []byte(goModuleNoTask),
					"/some/path/root/mod2/module.toml": []byte(goModuleWithTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := exec(dependencies, &flags.CommandFlags{
				Target: "all",
			}, &flags.GlobalFlags{}, "script_name")
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			expected := "the script is missing in modules mod1. Use --allow-missing or --exclude"
			if err.Error() != expected {
				t.Fatalf("Expected error message to be %s, got %s", expected, err.Error())
			}
		})
		t.Run("should execute all tasks if they are all present", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root/mod1",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootNoTask),
					"/some/path/root/mod1/module.toml": []byte(goModuleWithTask),
					"/some/path/root/mod2/module.toml": []byte(goModuleWithTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := exec(dependencies, &flags.CommandFlags{
				Target: "all",
			}, &flags.GlobalFlags{}, "script_name")
			if err != nil {
				t.Fatal(err)
			}
			commands := effects.Executor.Output()
			expected := "bash " + "some_bash_script"
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
	})
	t.Run("modules target with --allow-missing", func(t *testing.T) {
		t.Run("should execute all tasks if some are present", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root/mod1",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootNoTask),
					"/some/path/root/mod1/module.toml": []byte(goModuleNoTask),
					"/some/path/root/mod2/module.toml": []byte(goModuleWithTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := exec(dependencies, &flags.CommandFlags{
				Target:       "all",
				AllowMissing: true,
			}, &flags.GlobalFlags{}, "script_name")
			if err != nil {
				t.Fatal(err)
			}
			commands := effects.Executor.Output()
			expected := "bash " + "some_bash_script"
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
		t.Run("should fail if the task is missing in all modules", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root/mod1",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootNoTask),
					"/some/path/root/mod1/module.toml": []byte(goModuleNoTask),
					"/some/path/root/mod2/module.toml": []byte(goModuleNoTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := exec(dependencies, &flags.CommandFlags{
				Target:       "all",
				AllowMissing: true,
			}, &flags.GlobalFlags{}, "script_name")
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			expected := "the script is missing in all modules"
			if err.Error() != expected {
				t.Fatalf("Expected error message to be %s, got %s", expected, err.Error())
			}
		})
	})
	t.Run("modules target with --exclude", func(t *testing.T) {
		t.Run("should not execute the excluded modules", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root/mod1",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootNoTask),
					"/some/path/root/mod1/module.toml": []byte(goModuleWithTask),
					"/some/path/root/mod2/module.toml": []byte(goModuleWithTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := exec(dependencies, &flags.CommandFlags{
				Target:  "all",
				Exclude: "mod1",
			}, &flags.GlobalFlags{}, "script_name")
			if err != nil {
				t.Fatal(err)
			}
			commands := effects.Executor.Output()
			expected := "bash " + "some_bash_script"
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
	})
	t.Run("execute task queues", func(t *testing.T) {
		t.Run("should execute queue tasks at root", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootWithTaskQueue),
					"/some/path/root/mod1/module.toml": []byte(goModuleNoTask),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := exec(dependencies, &flags.CommandFlags{
				Target: "root",
			}, &flags.GlobalFlags{}, "script_name")
			if err != nil {
				t.Fatal(err)
			}
			commands := effects.Executor.Output()
			expected := "bash " + "some_bash_script"
			if len(commands) != 2 {
				t.Fatalf("Expected %d, got %d", 1, len(commands))
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
		t.Run("should execute queue tasks in a module", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(rootNoTask),
					"/some/path/root/mod1/module.toml": []byte(goModuleWithTaskQueue),
				},
			})
			cfg, _ := config.NewConfig(effects.ToEffects())
			dependencies := config.NewDependencies(effects.ToEffects(), cfg)
			err := exec(dependencies, &flags.CommandFlags{
				Target: "mod1",
			}, &flags.GlobalFlags{}, "script_name")
			if err != nil {
				t.Fatal(err)
			}
			commands := effects.Executor.Output()
			expected := "bash " + "some_bash_script"
			if len(commands) != 2 {
				t.Fatalf("Expected %d, got %d", 1, len(commands))
			}
			if commands[0].Dir != "/some/path/root/mod1" {
				t.Fatalf("Expected %s, got %s", "/some/path/root/mod1", commands[0].Dir)
			}
			if commands[0].Command != expected+"_1" {
				t.Fatalf("Expected %s, got %s", expected+"_1", commands[0].Command)
			}
			if commands[1].Dir != "/some/path/root/mod1" {
				t.Fatalf("Expected %s, got %s", "/some/path/root/mod1", commands[1].Dir)
			}
			if commands[1].Command != expected+"_2" {
				t.Fatalf("Expected %s, got %s", expected+"_2", commands[1].Command)
			}
		})
	})
}

var rootNoTask, _ = toml.Marshal(config.RootConfig{
	Name: "test",
})

var rootWithTask, _ = toml.Marshal(config.RootConfig{
	Name: "test",
	Tasks: map[string]config.TaskQueue{
		"script_name": config.TaskQueue{"some_bash_script"},
	},
})

var rootWithTaskQueue, _ = toml.Marshal(config.RootConfig{
	Name: "test",
	Tasks: map[string]config.TaskQueue{
		"script_name": config.TaskQueue{"some_bash_script_1", "some_bash_script_2"},
	},
})

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
		"script_name": config.TaskQueue{"some_bash_script"},
	},
})

var goModuleWithTaskQueue, _ = toml.Marshal(config.ModuleConfig{
	Template: "@default",
	Type:     "any",
	Language: "go",
	Tasks: map[string]config.TaskQueue{
		"script_name": config.TaskQueue{"some_bash_script_1", "some_bash_script_2"},
	},
})
