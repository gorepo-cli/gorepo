package initcmd

import (
	"github.com/pelletier/go-toml/v2"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"gorepo-cli/pkg"
	"testing"
)

func TestCommandInit(t *testing.T) {
	t.Run("should return an error if work.toml already exists", func(t *testing.T) {
		rootConfigBytes, _ := toml.Marshal(config.RootConfig{
			Name: "my-monorepo",
		})
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml": rootConfigBytes,
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := initCmd(dependencies, nil, nil, "")
		if err.Error() != "monorepo already exists at /some/path/root" {
			t.Fatalf("expected 'work.toml already exists at root', got %s", err.Error())
		}
	})
	t.Run("should create work.toml if there is no such file at root", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD:    "/some/path/root",
			Files: map[string][]byte{},
			QaBool: map[string]bool{
				"Do you want to vendor dependencies?": true,
			},
			QaString: map[string]string{
				"What is the monorepo name?": "",
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		_ = initCmd(dependencies, nil, &flags.GlobalFlags{
			Verbose: false,
		}, "")
		files := effects.Filesystem.Output()
		if files["/some/path/root/work.toml"] == nil {
			t.Fatal("expected a non-nil value, got nil")
		}
	})
	t.Run("should create a go.work file if there are none", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD:    "/some/path/root",
			Files: map[string][]byte{},
			QaBool: map[string]bool{
				"Do you want to vendor dependencies?": true,
			},
			QaString: map[string]string{
				"What is the monorepo name?": "test-monorepo",
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := initCmd(dependencies, nil, &flags.GlobalFlags{
			Verbose: false,
		}, "")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		executorOutput := effects.Executor.Output()
		found := false
		for _, cmd := range executorOutput {
			if cmd.Command == "go work init" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected 'go work init' to be called but it was not")
		}
	})
	t.Run("should not create a go.work file if there is one", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml": []byte{50},
			},
			QaBool: map[string]bool{
				"Do you want to vendor dependencies?": true,
			},
			QaString: map[string]string{
				"What is the monorepo name?": "test-monorepo",
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := initCmd(dependencies, nil, &flags.GlobalFlags{
			Verbose: false,
		}, "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		executorOutput := effects.Executor.Output()
		found := false
		for _, cmd := range executorOutput {
			if cmd.Command == "go work init" {
				found = true
				break
			}
		}
		if found {
			t.Fatal("expected 'go work init' not to be called but it was")
		}
	})
	t.Run("should use the folder name if the user does not select one", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD:    "/some/path/root",
			Files: map[string][]byte{},
			QaBool: map[string]bool{
				"Do you want to vendor dependencies?": true,
			},
			QaString: map[string]string{
				"What is the monorepo name?": "",
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := initCmd(dependencies, nil, &flags.GlobalFlags{
			Verbose: false,
		}, "")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		files := effects.Filesystem.Output()
		var rootConfig1 config.RootConfig
		err = toml.Unmarshal(files["/some/path/root/work.toml"], &rootConfig1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if rootConfig1.Name != "root" {
			t.Fatalf("expected rootConfig.Name to be 'root', got '%s'", rootConfig1.Name)
		}
	})
	t.Run("should use the name passed as an argument", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD:    "/some/path/root",
			Files: map[string][]byte{},
			QaBool: map[string]bool{
				"Do you want to vendor dependencies?": true,
			},
			QaString: map[string]string{
				"What is the monorepo name?": "",
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := initCmd(dependencies, nil, &flags.GlobalFlags{
			Verbose: false,
		}, "foo")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		files := effects.Filesystem.Output()
		var rootConfig1 config.RootConfig
		err = toml.Unmarshal(files["/some/path/root/work.toml"], &rootConfig1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if rootConfig1.Name != "foo" {
			t.Fatalf("expected rootConfig.Name to be 'foo', got '%s'", rootConfig1.Name)
		}
	})
	t.Run("should use the name provided by the user if they answer the question about it", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD:    "/some/path/root",
			Files: map[string][]byte{},
			QaBool: map[string]bool{
				"Do you want to vendor dependencies?": true,
			},
			QaString: map[string]string{
				"What is the monorepo name?": "bar",
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := initCmd(dependencies, nil, &flags.GlobalFlags{
			Verbose: false,
		}, "")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		files := effects.Filesystem.Output()
		var rootConfig1 config.RootConfig
		err = toml.Unmarshal(files["/some/path/root/work.toml"], &rootConfig1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if rootConfig1.Name != "bar" {
			t.Fatalf("expected rootConfig.Name to be 'bar', got '%s'", rootConfig1.Name)
		}
	})
	t.Run("should set vendoring to true if the user selected it", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD:    "/some/path/root",
			Files: map[string][]byte{},
			QaBool: map[string]bool{
				"Do you want to vendor dependencies?": true,
			},
			QaString: map[string]string{
				"What is the monorepo name?": "",
			},
		})
		cfg, _ := config.NewConfig(effects.ToEffects())
		dependencies := config.NewDependencies(effects.ToEffects(), cfg)
		err := initCmd(dependencies, nil, &flags.GlobalFlags{
			Verbose: false,
		}, "")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		files := effects.Filesystem.Output()
		var rootConfig1 config.RootConfig
		err = toml.Unmarshal(files["/some/path/root/work.toml"], &rootConfig1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if rootConfig1.Vendor != true {
			t.Fatalf("expected true got false")
		}
	})
}
