package concourse

import (
	"os"
	"os/exec"

	"code.cloudfoundry.org/commandrunner"
)

type Reconfigurer struct {
	commandRunner commandrunner.CommandRunner
}

func NewReconfigurer(commandRunner commandrunner.CommandRunner) *Reconfigurer {
	return &Reconfigurer{
		commandRunner: commandRunner,
	}
}

func (r *Reconfigurer) Reconfigure(target, pipeline, configPath, variablesPath string) error {
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
