package config

type TaskQueue []string

type RootConfig struct {
	Name    string               `toml:"name"`
	Version string               `toml:"version"`
	Vendor  bool                 `toml:"vendor"`
	Tasks   map[string]TaskQueue `toml:"tasks"`
}

// this type is needed to parse tasks, that can be strings or array of strings in the toml
type rootConfigRaw struct {
	RootConfig
	Tasks map[string]any `toml:"tasks"`
}

type ModuleConfig struct {
	Name         string               `toml:"-"`        // defined at runtime
	PathFromRoot string               `toml:"-"`        // defined at runtime
	Template     string               `toml:"template"` // default is @default
	Type         string               `toml:"type"`     // executable, library, script, static
	Language     string               `toml:"language"` // go, javascript...
	Main         string               `toml:"main"`     // entry point
	Priority     int                  `toml:"priority"` // build priority (higher comes first)
	Tasks        map[string]TaskQueue `toml:"tasks"`
}

// this type is needed to parse tasks, that can be strings or array of strings in the toml
type moduleConfigRaw struct {
	ModuleConfig
	Tasks map[string]any `toml:"tasks"`
}
