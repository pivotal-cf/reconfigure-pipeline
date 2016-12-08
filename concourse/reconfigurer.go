package concourse

import (
	"os"
	"os/exec"

	"code.cloudfoundry.org/commandrunner"
)

type Reconfigurer interface {
	Reconfigure(target, pipeline, configPath, variablesPath string) error
}

type reconfigurer struct {
	commandRunner commandrunner.CommandRunner
}

func NewReconfigurer(commandRunner commandrunner.CommandRunner) Reconfigurer {
	return &reconfigurer{
		commandRunner: commandRunner,
	}
}

func (r *reconfigurer) Reconfigure(target, pipeline, configPath, variablesPath string) error {
	args := []string{"-t", target, "set-pipeline", "-p", pipeline, "-c", configPath}
	if variablesPath != "" {
		args = append(args, "-l", variablesPath)
	}

	cmd := exec.Command("fly", args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return r.commandRunner.Run(cmd)
}
