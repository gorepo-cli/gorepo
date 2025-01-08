package flags

import "github.com/urfave/cli/v2"

var (
	Target = &cli.StringFlag{
		Name:  "target",
		Value: "root",
		Usage: "Target specific modules or root (comma separated)",
	}
	Exclude = &cli.StringFlag{
		Name:  "exclude",
		Value: "",
		Usage: "Exclude specific modules (comma separated)",
	}
	AllowMissing = &cli.BoolFlag{
		Name:  "allow-missing",
		Value: false,
		Usage: "Allow executing the scripts, even if some module don't have it",
	}
	Ci = &cli.BoolFlag{
		Name:  "ci",
		Value: false,
		Usage: "Enable mode CI",
	}
)

var (
	Verbose = &cli.BoolFlag{
		Name:  "verbose",
		Usage: "Enable verbose logging for all commands",
		Value: false,
	}
	Experimental = &cli.BoolFlag{
		Name:  "experimental",
		Value: false,
		Usage: "Experiment coming features",
	}
)

var ExecutionGroup = []cli.Flag{Target, Exclude}

var GlobalGroup = []cli.Flag{Verbose, Experimental}
