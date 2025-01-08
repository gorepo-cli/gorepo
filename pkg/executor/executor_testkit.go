package executor

import "strings"

type MockCommand struct {
	Dir     string
	Command string
	// We can improve this by enabling setting of the output and errors
	// Output  string
	// Err     error
}

type MockMethods interface {
	Methods
	Output() []MockCommand
}

type MockExecutor struct {
	Commands []MockCommand
}

var _ MockMethods = &MockExecutor{}

func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		Commands: []MockCommand{},
	}
}

func (m *MockExecutor) Go(absolutePath string, args ...string) error {
	cmd := strings.Join(args, " ")
	m.Commands = append(m.Commands, MockCommand{
		Dir:     absolutePath,
		Command: "go " + cmd,
	})
	return nil
}

func (m *MockExecutor) Bash(absolutePath, script string) error {
	m.Commands = append(m.Commands, MockCommand{
		Dir:     absolutePath,
		Command: script,
	})
	return nil
}

func (m *MockExecutor) Output() []MockCommand {
	return m.Commands
}
