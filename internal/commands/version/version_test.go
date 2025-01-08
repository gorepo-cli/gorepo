package version

import (
	"github.com/pelletier/go-toml/v2"
	"gorepo-cli/internal/config"
	"gorepo-cli/pkg"
	"testing"
)

func TestCommandVersion(t *testing.T) {
	t.Run("should log the version", func(t *testing.T) {
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
		_ = version(dependencies, nil, nil)
		logs := effects.Logger.Output()
		if logs[0] != "DEFAULT: dev" {
			t.Fatalf("Expected %s, got %s", "dev", logs[0])
		}
	})
}
