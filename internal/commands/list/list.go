package list

import (
	"bytes"
	"fmt"
	"gorepo-cli/internal/config"
	"gorepo-cli/internal/flags"
	"text/tabwriter"
)

func list(dependencies *config.Dependencies, cmdFlags *flags.CommandFlags, globalFlags *flags.GlobalFlags) error {
	if err := dependencies.Config.BreakIfRootConfigDoesNotExist(); err != nil {
		return err
	}
	modules, err := dependencies.Config.GetModules([]string{"all"}, []string{})
	if err != nil {
		return err
	}
	if len(modules) == 0 {
		dependencies.Effects.Logger.InfoLn("no modules found")
	} else {
		var buf bytes.Buffer
		w := tabwriter.NewWriter(&buf, 1, 1, 6, ' ', 0)
		for _, module := range modules {
			_, _ = fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%d\t%s\t%s\t%s\t", module.Name, module.RelativePath, module.Priority, module.Type, module.Template, module.Language))
		}
		_ = w.Flush()
		for _, line := range bytes.Split(buf.Bytes(), []byte("\n")) {
			if len(line) > 0 {
				dependencies.Effects.Logger.DefaultLn(string(line))
			}
		}
	}
	return nil
}
