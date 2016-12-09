package commandrunner

import (
	"os"
	"os/exec"
)

type SimpleCommandRunner struct{}

func NewSimpleCommandRunner() *SimpleCommandRunner {
	return &SimpleCommandRunner{}
}

func (r *SimpleCommandRunner) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}

func (r *SimpleCommandRunner) Start(cmd *exec.Cmd) error {
	return cmd.Start()
}

func (r *SimpleCommandRunner) Background(cmd *exec.Cmd) error {
	panic("SimpleCommandRunner does not support Background")
}

func (r *SimpleCommandRunner) Wait(cmd *exec.Cmd) error {
	return cmd.Wait()
}

func (r *SimpleCommandRunner) Kill(cmd *exec.Cmd) error {
	panic("SimpleCommandRunner does not support Kill")
	return nil
}

func (r *SimpleCommandRunner) Signal(cmd *exec.Cmd, signal os.Signal) error {
	panic("SimpleCommandRunner does not support Signal")
	return nil
}
