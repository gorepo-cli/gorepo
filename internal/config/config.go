package config

import (
	"errors"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"gorepo-cli/pkg"
	"os"
	"path/filepath"
	"sort"
)

type Config struct {
	Static  StaticConfig
	Runtime RuntimeConfig
	Effects *pkg.Effects
}

type StaticConfig struct {
	MaxRecursion   int
	RootFileName   string
	ModuleFileName string
}

type RuntimeConfig struct {
	WD   string
	ROOT string
	// BIN  string
	ForbiddenNames []string
}

type RootMethods interface {
	RootConfigExists() bool
	BreakIfRootConfigDoesNotExist() error
	GoWorkspaceExists() bool
	GetRootConfig() (cfg RootConfig, err error)
	WriteRootConfig(rootConfig RootConfig) (err error)
}

var _ RootMethods = &Config{}

type ModuleMethods interface {
	GetModules(targets, exclude []string) (modules []ModuleConfig, err error)
	FilterModulesByType(modules []ModuleConfig, types []string) []ModuleConfig
	FilterModulesByLanguage(modules []ModuleConfig, languages []string) []ModuleConfig
	GetModuleConfig(relativePath string) (cfg ModuleConfig, err error)
	WriteModuleConfig(modConfig ModuleConfig, absolutePathAndName string) (err error)
}

var _ ModuleMethods = &Config{}

type OtherMethods interface {
	PushForbiddenNames(names []string)
}

var _ OtherMethods = &Config{}

func NewConfig(effects *pkg.Effects) (cfg *Config, err error) {
	cfg = &Config{}
	cfg.Static = StaticConfig{
		MaxRecursion:   7,
		RootFileName:   "work.toml",
		ModuleFileName: "module.toml",
	}
	cfg.Runtime = RuntimeConfig{
		ForbiddenNames: []string{},
	}
	cfg.Effects = effects
	if wd, err := effects.Filesystem.GetWd(); err == nil {
		cfg.Runtime.WD = wd
	} else {
		return cfg, err
	}
	if root, err := getRootPath(cfg); err == nil {
		cfg.Runtime.ROOT = root
	} else {
		return cfg, err
	}
	return cfg, nil
}

func getRootPath(cfg *Config) (root string, err error) {
	currentDir := cfg.Runtime.WD
	if currentDir == "" {
		return "", fmt.Errorf("no working directory")
	}
	for i := 0; i <= cfg.Static.MaxRecursion; i++ {
		filePath := filepath.Join(currentDir, cfg.Static.RootFileName)
		if exists := cfg.Effects.Filesystem.Exists(filePath); exists {
			return currentDir, nil
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return cfg.Runtime.WD, nil
		}
		currentDir = parentDir
	}
	return "", fmt.Errorf("root not found")
}

func (c *Config) RootConfigExists() bool {
	filePath := filepath.Join(c.Runtime.ROOT, c.Static.RootFileName)
	return c.Effects.Filesystem.Exists(filePath)
}

func (c *Config) BreakIfRootConfigDoesNotExist() error {
	if !c.RootConfigExists() {
		return fmt.Errorf("monorepo configuration not found at %s", c.Runtime.ROOT)
	}
	return nil
}

func (c *Config) GetRootConfig() (cfg RootConfig, err error) {
	cfgRaw := &rootConfigRaw{}
	file, err := c.Effects.Filesystem.Read(filepath.Join(c.Runtime.ROOT, c.Static.RootFileName))
	if err != nil {
		return cfg, err
	}
	err = toml.Unmarshal(file, &cfgRaw)
	if err != nil {
		return cfg, err
	}
	// todo: duplicating everything one by one is naive
	cfg.Name = cfgRaw.Name
	cfg.Version = cfgRaw.Version
	cfg.Vendor = cfgRaw.Vendor
	for scriptName, script := range cfgRaw.Tasks {
		if _, ok := script.(string); ok {
			if cfg.Tasks == nil {
				cfg.Tasks = make(map[string]TaskQueue)
			}
			cfg.Tasks[scriptName] = []string{script.(string)}
		} else if _, ok := script.([]interface{}); ok {
			if cfg.Tasks == nil {
				cfg.Tasks = make(map[string]TaskQueue)
			}
			var scriptQueue []string
			for _, item := range script.([]interface{}) {
				if str, ok := item.(string); ok {
					scriptQueue = append(scriptQueue, str)
				} else {
					return cfg, fmt.Errorf("invalid script format in %s: non-string value found", scriptName)
				}
			}
			cfg.Tasks[scriptName] = scriptQueue
		}
	}
	return cfg, nil
}

func (c *Config) WriteRootConfig(rootConfig RootConfig) (err error) {
	configStr, err := toml.Marshal(rootConfig)
	if err != nil {
		return err
	}
	filePath := filepath.Join(c.Runtime.ROOT, c.Static.RootFileName)
	return c.Effects.Filesystem.Write(filePath, configStr)
}

func (c *Config) GoWorkspaceExists() bool {
	filePath := filepath.Join(c.Runtime.ROOT, "go.work")
	return c.Effects.Filesystem.Exists(filePath)
}

func (c *Config) GetModules(targets, exclude []string) (modules []ModuleConfig, err error) {
	// validation
	for _, target := range targets {
		if target == "root" && len(targets) > 1 {
			return nil, errors.New("cannot run script in root and in modules at the same time, you're being too greedy, run the command twice")
		} else if target == "all" && len(targets) > 1 {
			return nil, errors.New("cannot run script in all modules and in specific modules, non sense")
		}
	}
	for _, excluded := range exclude {
		if excluded == "all" {
			return nil, errors.New("excluding all modules makes no sense")
		} else if excluded == "root" {
			return nil, errors.New("excluding root is the default behaviour, no need to specify it")
		}
	}
	// walk
	currentPath := c.Runtime.ROOT
	skipDirs := map[string]bool{
		"node_modules": true,
		"vendor":       true,
		"bin":          true,
		"obj":          true,
		"target":       true,
		".cache":       true,
		".venv":        true,
		"virtualenvs":  true,
		"venv":         true,
		"env":          true,
		".git":         true,
		".github":      true,
		".gitlab":      true,
		".vscode":      true,
		"__pycache__":  true,
		"dist":         true,
		"out":          true,
		".snv":         true,
		".hg":          true,
		".DS_Store":    true,
		"thumbs.db":    true,
		"logs":         true,
		".idea":        true,
		".settings":    true,
		".terraform":   true,
		".aws-sam":     true,
		".circleci":    true,
	}
	err = c.Effects.Filesystem.Walk(currentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && skipDirs[info.Name()] {
			return filepath.SkipDir
		}
		if info.IsDir() {
			exists := c.Effects.Filesystem.Exists(filepath.Join(path, c.Static.ModuleFileName))
			if exists {
				relativePath, err := filepath.Rel(c.Runtime.ROOT, path)
				if err != nil {
					return err
				}
				moduleConfig, err := c.GetModuleConfig(relativePath)
				if err != nil {
					return err
				}
				if (targets[0] == "all" || contains(targets, moduleConfig.Name)) && !contains(exclude, moduleConfig.Name) {
					modules = append(modules, moduleConfig)
				}
			}
		}
		return nil
	})
	if err != nil {
		c.Effects.Logger.WarningLn(err.Error())
		return modules, err
	}
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Name < modules[j].Name
	})
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Priority > modules[j].Priority
	})
	return modules, nil
}

func (c *Config) FilterModulesByType(modules []ModuleConfig, types []string) []ModuleConfig {
	var filteredModules []ModuleConfig
	for _, module := range modules {
		if contains(types, module.Type) {
			filteredModules = append(filteredModules, module)
		}
	}
	return filteredModules
}

func (c *Config) FilterModulesByLanguage(modules []ModuleConfig, languages []string) []ModuleConfig {
	var filteredModules []ModuleConfig
	for _, module := range modules {
		if contains(languages, module.Language) {
			filteredModules = append(filteredModules, module)
		}
	}
	return filteredModules
}

func (c *Config) GetModuleConfig(relativePath string) (cfg ModuleConfig, err error) {
	cfgRaw := &moduleConfigRaw{}
	path := filepath.Join(c.Runtime.ROOT, relativePath, c.Static.ModuleFileName)
	file, err := c.Effects.Filesystem.Read(path)
	if err != nil {
		return cfg, err
	}
	err = toml.Unmarshal(file, &cfgRaw)
	if err != nil {
		return cfg, err
	}
	// todo: duplicating everything one by one is naive
	cfg.Name = cfgRaw.Name
	cfg.PathFromRoot = cfgRaw.PathFromRoot
	cfg.Template = cfgRaw.Template
	cfg.Type = cfgRaw.Type
	cfg.Language = cfgRaw.Language
	cfg.Main = cfgRaw.Main
	cfg.Priority = cfgRaw.Priority
	for scriptName, script := range cfgRaw.Tasks {
		if _, ok := script.(string); ok {
			if cfg.Tasks == nil {
				cfg.Tasks = make(map[string]TaskQueue)
			}
			cfg.Tasks[scriptName] = []string{script.(string)}
		} else if _, ok := script.([]interface{}); ok {
			if cfg.Tasks == nil {
				cfg.Tasks = make(map[string]TaskQueue)
			}
			var scriptQueue []string
			for _, item := range script.([]interface{}) {
				if str, ok := item.(string); ok {
					scriptQueue = append(scriptQueue, str)
				} else {
					return cfg, fmt.Errorf("invalid script format in %s: non-string value found", scriptName)
				}
			}
			cfg.Tasks[scriptName] = scriptQueue
		}
	}
	cfg.Name = filepath.Base(relativePath)
	cfg.PathFromRoot = relativePath
	return cfg, nil
}

func (c *Config) WriteModuleConfig(modConfig ModuleConfig, absolutePathAndName string) (err error) {
	configStr, err := toml.Marshal(modConfig)
	if err != nil {
		return err
	}
	err = c.Effects.Filesystem.MkDir(absolutePathAndName)
	if err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	filePath := filepath.Join(absolutePathAndName, c.Static.ModuleFileName)
	return c.Effects.Filesystem.Write(filePath, configStr)
}

func (c *Config) PushForbiddenNames(names []string) {
	for _, name := range names {
		c.Runtime.ForbiddenNames = append(c.Runtime.ForbiddenNames, name)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
