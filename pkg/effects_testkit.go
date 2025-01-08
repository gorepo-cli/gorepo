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

func NewTestkit(args TestKitArgs) (tk *TestKitResponse) {
	var (
		_executor   = executor.NewMockExecutor()
		_filesystem = filesystem.NewMockFilesystem(args.Files, args.WD)
		_logger     = logger.NewMockLogger()
		_terminal   = terminal.NewMockTerminal(args.WD, args.QaBool, args.QaString)
	)
	return &TestKitResponse{
		Effects: &Effects{
			Executor:   _executor,
			Filesystem: _filesystem,
			Logger:     _logger,
			Terminal:   _terminal,
		},
		GetExecutorOutput:   _executor.Output,
		GetFilesystemOutput: _filesystem.Output,
		GetLoggerOutput:     _logger.Output,
	}
}
