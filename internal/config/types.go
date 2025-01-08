package config

type RootConfig struct {
	Name     string            `toml:"name"`
	Version  string            `toml:"version"`
	Strategy string            `toml:"strategy"` // workspace / rewrites (unsupported)
	Vendor   bool              `toml:"vendor"`   // vendor or not (unsupported)
	Scripts  map[string]string `toml:"scripts"`
}

// ModuleConfig contains the configuration of a module
type ModuleConfig struct {
	// Module's name (= folder's name), added at runtime
	Name string `toml:"-"`
	// Relative path to the root, added at runtime
	RelativePath string `toml:"-"`
	// Name of the template (default is @default)
	Template string `toml:"template"`
	// Module's type (executable (can be built and executed), library (can be built), static (can not be built))
	Type string `toml:"type"`
	// Language of the module (go, python, node, javascript, etc.)
	Language string `toml:"language"`
	// Entry point of the module, if needed to be built
	Main string `toml:"main"`
	// Build priority, higher goes first
	Priority int `toml:"priority"`
	// List of scripts that can be run through gorepo execute <script_name>
	Scripts map[string]string `toml:"scripts"`
}
