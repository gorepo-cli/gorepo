package config

import (
	"gorepo-cli/pkg"
	"testing"
)

var tomlWithStringScripts = `
name = "test"
[scripts]
test = "echo string_script"
`

var tomlWithArrayScripts = `
name = "test"
[scripts]
test = ["echo array_script_1", "echo array_script_2"]
`

var tomlWithoutScripts = `
name = "test"
`

func TestConfig_GetRootConfig(t *testing.T) {
	t.Run("should parse scripts that are a single strings", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml": []byte(tomlWithStringScripts),
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
		//if len(rootCfg.Scripts["test"]) != 1 {
		//	t.Fatalf("Expected %d, got %d", 1, len(rootCfg.Scripts["test"]))
		//}
		if rootCfg.Scripts["test"] != "echo string_script" {
			t.Fatalf("Expected %s, got %s", "echo string_script", rootCfg.Scripts["test"])
		}
	})
	//t.Run("should parse scripts that are an array of strings", func(t *testing.T) {
	//	effects := pkg.NewTestkit(pkg.TestKitArgs{
	//		WD: "/some/path/root",
	//		Files: map[string][]byte{
	//			"/some/path/root/work.toml": []byte(tomlWithArrayScripts),
	//		},
	//	})
	//	cfg, err := NewConfig(effects.ToEffects())
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	rootCfg, _ := cfg.GetRootConfig()
	//	if rootCfg.Name != "test" {
	//		t.Fatalf("Expected %s, got %s", "test", rootCfg.Name)
	//	}
	//	if len(rootCfg.Scripts["test"]) != 2 {
	//		t.Fatalf("Expected %d, got %d", 1, len(rootCfg.Scripts["test"]))
	//	}
	//	if rootCfg.Scripts["test"][0] != "echo array_script_1" {
	//		t.Fatalf("Expected %s, got %s", "echo array_script_1", rootCfg.Scripts["test"][0])
	//	}
	//	if rootCfg.Scripts["test"][1] != "echo array_script_2" {
	//		t.Fatalf("Expected %s, got %s", "echo array_script_2", rootCfg.Scripts["test"][1])
	//	}
	//})
	t.Run("should not break if there are no scripts", func(t *testing.T) {
		effects := pkg.NewTestkit(pkg.TestKitArgs{
			WD: "/some/path/root",
			Files: map[string][]byte{
				"/some/path/root/work.toml": []byte(tomlWithoutScripts),
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
		if len(rootCfg.Scripts["test"]) != 0 {
			t.Fatalf("Expected %d, got %d", 0, len(rootCfg.Scripts["test"]))
		}
	})
}
