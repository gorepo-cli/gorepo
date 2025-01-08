package main

import (
	"gorepo-cli/internal/execute"
	"gorepo-cli/pkg/logger"
	"os"
)

func main() {
	l := logger.NewLogger()
	if err := execute.Execute(l); err != nil {
		l.FatalLn(err.Error())
		os.Exit(1)
	}
}
