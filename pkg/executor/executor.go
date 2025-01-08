package executor

import (
	"fmt"
	"os"
	"os/exec"
)

type Methods interface {
	Go(absolutePath string, args ...string) error
	Bash(absolutePath, script string) error
}

type Executor struct{}

var _ Methods = &Executor{}

func (x *Executor) Go(absolutePath string, args ...string) (err error) {
	cmd := exec.Command("go", args...)
	cmd.Dir = absolutePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run command: %w\nOutput: %s", err, string(output))
	}
	return nil
}

func (x *Executor) Bash(absolutePath, script string) (err error) {
	cmd := exec.Command("/bin/sh", "-c", script)
	cmd.Dir = absolutePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command in %s: %w", absolutePath, err)
	}
	return nil
}
