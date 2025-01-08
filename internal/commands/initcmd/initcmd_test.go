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
	t.Run("should create a go.work file if it is missing", func(t *testing.T) {
		// todo
	})
	t.Run("should use the folder name if the user does not select one", func(t *testing.T) {
		// todo
	})
	t.Run("should use the name passed as an argument", func(t *testing.T) {
		// todo
	})
	t.Run("should use the name provided by the user if they answer the question about it", func(t *testing.T) {
		// todo
	})
	t.Run("should set vendoring to true if the user selected it", func(t *testing.T) {
		// todo
	})
}
