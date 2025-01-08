package pkg

import (
	"gorepo-cli/pkg/executor"
	"gorepo-cli/pkg/filesystem"
	"gorepo-cli/pkg/logger"
	"gorepo-cli/pkg/terminal"
)

type Effects struct {
	Executor   executor.Methods
	Filesystem filesystem.Methods
	Logger     logger.Methods
	Terminal   terminal.Methods
}

func NewEffects(l logger.Methods) *Effects {
	return &Effects{
		Executor:   &executor.Executor{},
		Filesystem: &filesystem.Filesystem{},
		Logger:     l,
		Terminal:   &terminal.Terminal{},
	}
}
