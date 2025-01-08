package flags

import "github.com/urfave/cli/v2"

type CommandFlags struct {
	Target       string
	Exclude      string
	AllowMissing bool
	Ci           bool
}

func ExtractCommandFlags(c *cli.Context) (executionFlags *CommandFlags) {
	return &CommandFlags{
		Target:       c.String(Target.Name),
		Exclude:      c.String(Exclude.Name),
		AllowMissing: c.Bool(AllowMissing.Name),
		Ci:           c.Bool(Ci.Name),
	}
}

type GlobalFlags struct {
	Verbose      bool
	Experimental bool
}

func ExtractGlobalFlags(c *cli.Context) (globalFlags *GlobalFlags) {
	return &GlobalFlags{
		Verbose:      c.Bool(Verbose.Name),
		Experimental: c.Bool(Experimental.Name),
	}
}
