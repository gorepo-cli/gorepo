package commands

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
	"gorepo-cli/internal/config"
	"gorepo-cli/pkg"
	"strconv"
	"testing"
)

// entities that could conflict:
// - command-name
// - module-name
// - module-subcommand-name
// - root-task-name
// - module-task-name

func TestRegisterCommands(t *testing.T) {
	t.Run("should register commands if no conflict exist", func(t *testing.T) {
		var rootConfig, _ = toml.Marshal(config.RootConfig{
			Name: "test",
			Tasks: map[string]config.TaskQueue{
				"script_name": config.TaskQueue{"some_bash_script_1"},
			},
		})
		var modConfig1, _ = toml.Marshal(config.ModuleConfig{
			Template: "@default",
			Tasks: map[string]config.TaskQueue{
				"script_name": config.TaskQueue{"some_bash_script_2"},
			},
		})
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml":        []byte(rootConfig),
				"/some/path/root/mod1/module.toml": []byte(modConfig1),
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		var commands []*cli.Command
		cmd, err := RegisterCommands(commands, func(dependencies *config.Dependencies) []*cli.Command {
			return []*cli.Command{
				{
					Name: "list",
				},
				{
					Name: "fmt",
				},
			}
		}, func(name string, dependencies *config.Dependencies) []*cli.Command {
			return []*cli.Command{
				{
					Name: "fmt",
				},
			}
		}, dependencies)
		if err != nil {
			t.Fatal(err)
		}
		if cmd[0].Name != "list" {
			t.Fatal("Expected list, got " + cmd[0].Name)
		}
		if cmd[1].Name != "fmt" {
			t.Fatal("Expected fmt, got " + cmd[1].Name)
		}
		if cmd[2].Name != "mod1" {
			t.Fatal("Expected mod1, got " + cmd[2].Name)
		}
		if cmd[2].Subcommands[0].Name != "fmt" {
			t.Fatal("Expected fmt, got " + cmd[2].Subcommands[0].Name)
		}
		if cmd[2].Subcommands[1].Name != "script_name" {
			t.Fatal("Expected script_name, got " + cmd[2].Subcommands[1].Name)
		}
		if cmd[3].Name != "script_name" {
			t.Fatal("Expected script_name, got " + cmd[3].Name)
		}
		if len(cmd) != 4 {
			t.Fatal("Expected 4 commands, got " + strconv.Itoa(len(cmd)))
		}
	})
	t.Run("should drop the module-name in case of a conflict between: command and module-name", func(t *testing.T) {
		// Example: if a module is named `exec` like the command, don't register the module
		var rootConfig, _ = toml.Marshal(config.RootConfig{
			Name: "test",
			Tasks: map[string]config.TaskQueue{
				"script_name": config.TaskQueue{"some_bash_script_1"},
			},
		})
		var modConfig1, _ = toml.Marshal(config.ModuleConfig{
			Template: "@default",
			Tasks: map[string]config.TaskQueue{
				"script_name": config.TaskQueue{"some_bash_script_2"},
			},
		})
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml":        []byte(rootConfig),
				"/some/path/root/list/module.toml": []byte(modConfig1),
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		var commands []*cli.Command
		cmd, err := RegisterCommands(commands, func(dependencies *config.Dependencies) []*cli.Command {
			return []*cli.Command{
				{
					Name:  "list",
					Usage: "root-command",
				},
				{
					Name:  "fmt",
					Usage: "root-command",
				},
			}
		}, func(name string, dependencies *config.Dependencies) []*cli.Command {
			return []*cli.Command{
				{
					Name:  "fmt",
					Usage: "module-sub-command",
				},
			}
		}, dependencies)
		if err != nil {
			t.Fatal(err)
		}
		if cmd[0].Name != "list" && cmd[0].Usage != "root-command" {
			t.Fatal("Expected list, got " + cmd[0].Name)
		}
		if cmd[1].Name != "fmt" {
			t.Fatal("Expected fmt, got " + cmd[1].Name)
		}
		if cmd[2].Name != "script_name" {
			t.Fatal("Expected script_name, got " + cmd[3].Name)
		}
		if len(cmd) != 3 {
			t.Fatal("Expected 3 commands, got " + strconv.Itoa(len(cmd)))
		}
	})
	t.Run("should drop the root-task in case of a conflict between: command and root-task", func(t *testing.T) {
		// Example: if a task is called fmt, like the command, don't register the command
		var rootConfig, _ = toml.Marshal(config.RootConfig{
			Name: "test",
			Tasks: map[string]config.TaskQueue{
				"list": config.TaskQueue{"some_bash_script_1"},
			},
		})
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml": []byte(rootConfig),
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		var commands []*cli.Command
		cmd, err := RegisterCommands(commands, func(dependencies *config.Dependencies) []*cli.Command {
			return []*cli.Command{
				{
					Name:  "list",
					Usage: "root-command",
				},
				{
					Name:  "fmt",
					Usage: "root-command",
				},
			}
		}, func(name string, dependencies *config.Dependencies) []*cli.Command {
			return []*cli.Command{
				{
					Name:  "fmt",
					Usage: "module-sub-command",
				},
			}
		}, dependencies)
		if err != nil {
			t.Fatal(err)
		}
		if cmd[0].Name != "list" && cmd[0].Usage != "root-command" {
			t.Fatal("Expected list, got " + cmd[0].Name)
		}
		if cmd[1].Name != "fmt" {
			t.Fatal("Expected fmt, got " + cmd[1].Name)
		}
		if len(cmd) != 2 {
			t.Fatal("Expected 2 commands, got " + strconv.Itoa(len(cmd)))
		}
	})
	t.Run("should drop the root task in case of a conflict between: module command and root task", func(t *testing.T) {
		// Example: if a module is called foo and a root-task is named foo: don't register the task
		var rootConfig, _ = toml.Marshal(config.RootConfig{
			Name: "test",
			Tasks: map[string]config.TaskQueue{
				"foo": config.TaskQueue{"some_bash_script_1"},
			},
		})
		var modConfig1, _ = toml.Marshal(config.ModuleConfig{
			Template: "@default",
			Tasks: map[string]config.TaskQueue{
				"script_name": config.TaskQueue{"some_bash_script_2"},
			},
		})
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml":       []byte(rootConfig),
				"/some/path/root/foo/module.toml": []byte(modConfig1),
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		var commands []*cli.Command
		cmd, err := RegisterCommands(commands, func(dependencies *config.Dependencies) []*cli.Command {
			return []*cli.Command{
				{
					Name:  "list",
					Usage: "root-command",
				},
				{
					Name:  "fmt",
					Usage: "root-command",
				},
			}
		}, func(name string, dependencies *config.Dependencies) []*cli.Command {
			return []*cli.Command{
				{
					Name:  "fmt",
					Usage: "module-sub-command",
				},
			}
		}, dependencies)
		if err != nil {
			t.Fatal(err)
		}
		if cmd[0].Name != "list" && cmd[0].Usage != "root-command" {
			t.Fatal("Expected list, got " + cmd[0].Name)
		}
		if cmd[1].Name != "fmt" {
			t.Fatal("Expected fmt, got " + cmd[1].Name)
		}
		if cmd[2].Name != "foo" {
			t.Fatal("Expected foo, got " + cmd[2].Name)
		}
		if cmd[2].Subcommands == nil {
			t.Fatal("Expected subcommands, got nil")
		}
		if len(cmd) != 3 {
			t.Fatal("Expected 3 commands, got " + strconv.Itoa(len(cmd)))
		}
	})
	t.Run("should allow same name for module subcommand and module task", func(t *testing.T) {
		// Example: A module `deploy` with a script named `deploy`
		var rootConfig, _ = toml.Marshal(config.RootConfig{
			Name: "test",
		})
		var modConfig1, _ = toml.Marshal(config.ModuleConfig{
			Template: "@default",
			Tasks: map[string]config.TaskQueue{
				"foo": config.TaskQueue{"some_bash_script_2"},
			},
		})
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml":       []byte(rootConfig),
				"/some/path/root/foo/module.toml": []byte(modConfig1),
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		var commands []*cli.Command
		cmd, err := RegisterCommands(commands, func(dependencies *config.Dependencies) []*cli.Command {
			return []*cli.Command{
				{
					Name:  "list",
					Usage: "root-command",
				},
				{
					Name:  "fmt",
					Usage: "root-command",
				},
			}
		}, func(name string, dependencies *config.Dependencies) []*cli.Command {
			return []*cli.Command{
				{
					Name:  "fmt",
					Usage: "module-sub-command",
				},
			}
		}, dependencies)
		if err != nil {
			t.Fatal(err)
		}
		if cmd[0].Name != "list" && cmd[0].Usage != "root-command" {
			t.Fatal("Expected list, got " + cmd[0].Name)
		}
		if cmd[1].Name != "fmt" {
			t.Fatal("Expected fmt, got " + cmd[1].Name)
		}
		if cmd[2].Name != "foo" {
			t.Fatal("Expected foo, got " + cmd[2].Name)
		}
		if cmd[2].Subcommands[1].Name != "foo" {
			t.Fatal("Expected foo, got " + cmd[2].Subcommands[1].Name)
		}
		if len(cmd) != 3 {
			t.Fatal("Expected 3 commands, got " + strconv.Itoa(len(cmd)))
		}
	})
	t.Run("should allow root task and module task to have the same name", func(t *testing.T) {})
	t.Run("should allow a module task to have a command name", func(t *testing.T) {})
}
