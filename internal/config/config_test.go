package config

import (
	"gorepo-cli/pkg"
	"testing"
)

func TestConfig_NewConfig(t *testing.T) {
	t.Run("should detect WD", func(t *testing.T) {})
	t.Run("should detect ROOT if there is no work.toml", func(t *testing.T) {})
	t.Run("should detect ROOT if there is work.toml", func(t *testing.T) {})
}

func TestConfig_RootConfigExists(t *testing.T) {
	t.Run("should return false if there is no work.toml", func(t *testing.T) {})
	t.Run("should return true if there is work.toml in the same folder", func(t *testing.T) {})
	t.Run("should return true if there is work.toml in a parent folder", func(t *testing.T) {})
}

func TestConfig_WriteRootConfig(t *testing.T) {
	t.Run("should save root config", func(t *testing.T) {})
}

func TestConfig_GoWorkspaceExists(t *testing.T) {
	t.Run("should return false if there is no go.work", func(t *testing.T) {})
	t.Run("should return true if there is a go.work in the same folder", func(t *testing.T) {})
	t.Run("should return true if there is a go.work in a parent folder", func(t *testing.T) {})
}

var tomlWithStringTasks = `
name = "test"
[tasks]
test = "echo string_script"
`

var tomlWithArrayTasks = `
name = "test"
[tasks]
test = ["echo array_script_1", "echo array_script_2"]
`

var tomlWithoutTasks = `
name = "test"
`

func TestConfig_GetRootConfig(t *testing.T) {
	t.Run("tasks", func(t *testing.T) {
		t.Run("should parse tasks that are a single strings", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml": []byte(tomlWithStringTasks),
				},
			})
			cfg, err := NewConfig(effects.ToEffects())
			if err != nil {
				t.Fatal(err)
			}
			rootCfg, err := cfg.GetRootConfig()
			if err != nil {
				t.Fatal(err)
			}
			if rootCfg.Name != "test" {
				t.Fatalf("Expected %s, got %s", "test", rootCfg.Name)
			}
			if len(rootCfg.Tasks["test"]) != 1 {
				t.Fatalf("Expected %d, got %d", 1, len(rootCfg.Tasks["test"]))
			}
			if rootCfg.Tasks["test"][0] != "echo string_script" {
				t.Fatalf("Expected %s, got %s", "echo string_script", rootCfg.Tasks["test"])
			}
		})
		t.Run("should parse tasks that are an array of strings", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml": []byte(tomlWithArrayTasks),
				},
			})
			cfg, err := NewConfig(effects.ToEffects())
			if err != nil {
				t.Fatal(err)
			}
			rootCfg, _ := cfg.GetRootConfig()
			if rootCfg.Name != "test" {
				t.Fatalf("Expected %s, got %s", "test", rootCfg.Name)
			}
			if len(rootCfg.Tasks["test"]) != 2 {
				t.Fatalf("Expected %d, got %d", 2, len(rootCfg.Tasks["test"]))
			}
			if rootCfg.Tasks["test"][0] != "echo array_script_1" {
				t.Fatalf("Expected %s, got %s", "echo array_script_1", rootCfg.Tasks["test"][0])
			}
			if rootCfg.Tasks["test"][1] != "echo array_script_2" {
				t.Fatalf("Expected %s, got %s", "echo array_script_2", rootCfg.Tasks["test"][1])
			}
		})
		t.Run("should not break if there are no tasks", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml": []byte(tomlWithoutTasks),
				},
			})
			cfg, err := NewConfig(effects.ToEffects())
			if err != nil {
				t.Fatal(err)
			}
			rootCfg, _ := cfg.GetRootConfig()
			if rootCfg.Name != "test" {
				t.Fatalf("Expected %s, got %s", "test", rootCfg.Name)
			}
			if len(rootCfg.Tasks["test"]) != 0 {
				t.Fatalf("Expected %d, got %d", 0, len(rootCfg.Tasks["test"]))
			}
		})
	})
}

func TestConfig_GetModuleConfig(t *testing.T) {
	t.Run("tasks", func(t *testing.T) {
		t.Run("should parse tasks that are a single strings", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(tomlWithoutTasks),
					"/some/path/root/mod1/module.toml": []byte(tomlWithStringTasks),
				},
			})
			cfg, err := NewConfig(effects.ToEffects())
			if err != nil {
				t.Fatal(err)
			}
			modCfg, err := cfg.GetModuleConfig("mod1")
			if err != nil {
				t.Fatal(err)
			}
			if modCfg.Name != "mod1" {
				t.Fatalf("Expected %s, got %s", "mod1", modCfg.Name)
			}
			if len(modCfg.Tasks["test"]) != 1 {
				t.Fatalf("Expected %d, got %d", 1, len(modCfg.Tasks["test"]))
			}
			if modCfg.Tasks["test"][0] != "echo string_script" {
				t.Fatalf("Expected %s, got %s", "echo string_script", modCfg.Tasks["test"])
			}
		})
		t.Run("should parse tasks that are an array of strings", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml":        []byte(tomlWithArrayTasks),
					"/some/path/root/mod1/module.toml": []byte(tomlWithArrayTasks),
				},
			})
			cfg, err := NewConfig(effects.ToEffects())
			if err != nil {
				t.Fatal(err)
			}
			modCfg, _ := cfg.GetModuleConfig("mod1")
			if modCfg.Name != "mod1" {
				t.Fatalf("Expected %s, got %s", "mod1", modCfg.Name)
			}
			if len(modCfg.Tasks["test"]) != 2 {
				t.Fatalf("Expected %d, got %d", 2, len(modCfg.Tasks["test"]))
			}
			if modCfg.Tasks["test"][0] != "echo array_script_1" {
				t.Fatalf("Expected %s, got %s", "echo array_script_1", modCfg.Tasks["test"][0])
			}
			if modCfg.Tasks["test"][1] != "echo array_script_2" {
				t.Fatalf("Expected %s, got %s", "echo array_script_2", modCfg.Tasks["test"][1])
			}
		})
		t.Run("should not break if there are no tasks", func(t *testing.T) {
			effects := pkg.NewTestkit(pkg.TestKitArgs{
				WD: "/some/path/root",
				Files: map[string][]byte{
					"/some/path/root/work.toml": []byte(tomlWithoutTasks),
				},
			})
			cfg, err := NewConfig(effects.ToEffects())
			if err != nil {
				t.Fatal(err)
			}
			modCfg, _ := cfg.GetRootConfig()
			if modCfg.Name != "test" {
				t.Fatalf("Expected %s, got %s", "test", modCfg.Name)
			}
			if len(modCfg.Tasks["test"]) != 0 {
				t.Fatalf("Expected %d, got %d", 0, len(modCfg.Tasks["test"]))
			}
		})
	})
}

func TestConfig_GetModules(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		t.Run("should return no module if there is no module", func(t *testing.T) {})
		t.Run("should return modules, even nested", func(t *testing.T) {})
		t.Run("should return modules in alphabetical order (if no priority)", func(t *testing.T) {})
		t.Run("should return modules in priority order", func(t *testing.T) {})
	})
	t.Run("filtering", func(t *testing.T) {
		t.Run("should filter by target", func(t *testing.T) {})
		t.Run("should filter by exclude", func(t *testing.T) {})
	})
}

func TestConfig_FilterModulesByType(t *testing.T) {
	t.Run("should not filter if there is no type", func(t *testing.T) {})
	t.Run("should filter by type", func(t *testing.T) {})
}

func TestConfig_FilterModulesByLanguage(t *testing.T) {
	t.Run("should not filter if there is no language", func(t *testing.T) {})
	t.Run("should filter by language", func(t *testing.T) {})
}

func TestConfig_WriteModuleConfig(t *testing.T) {
	t.Run("should write the module config at path", func(t *testing.T) {})
}
