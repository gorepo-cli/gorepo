package pkg

import (
	"gorepo-cli/pkg/executor"
	"gorepo-cli/pkg/filesystem"
	"gorepo-cli/pkg/logger"
	"gorepo-cli/pkg/terminal"
)

type TestKitArgs struct {
	WD       string
	Files    map[string][]byte
	QaBool   map[string]bool
	QaString map[string]string
}

type TestKitResponse struct {
	Effects             *Effects
	GetFilesystemOutput func() map[string][]byte
	GetLoggerOutput     func() []string
	GetExecutorOutput   func() []executor.MockCommand
}

type MockEffects struct {
	Executor   executor.MockMethods
	Filesystem filesystem.MockMethods
	Logger     logger.MockMethods
	Terminal   terminal.MockMethods
}

// a bit naive
func (mock *MockEffects) ToEffects() *Effects {
	return &Effects{
		Executor:   mock.Executor,
		Filesystem: mock.Filesystem,
		Logger:     mock.Logger,
		Terminal:   mock.Terminal,
	}
}

func NewTestkit(args TestKitArgs) (effects *MockEffects) {
	var (
		_executor   = executor.NewMockExecutor()
		_filesystem = filesystem.NewMockFilesystem(args.Files, args.WD)
		_logger     = logger.NewMockLogger()
		_terminal   = terminal.NewMockTerminal(args.WD, args.QaBool, args.QaString)
	)
	return &MockEffects{
		Executor:   _executor,
		Filesystem: _filesystem,
		Logger:     _logger,
		Terminal:   _terminal,
	}
}
