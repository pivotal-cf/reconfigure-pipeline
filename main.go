package main

import (
	"flag"
	"os"

	"code.cloudfoundry.org/commandrunner/linux_command_runner"
	"github.com/oozie/reconfigure-pipeline/actions"
	"github.com/oozie/reconfigure-pipeline/concourse"
	"github.com/oozie/reconfigure-pipeline/fifo"
	"github.com/oozie/reconfigure-pipeline/lastpass"
)

func main() {
	var pipeline, configPath, target, variablesPath string

	flag.StringVar(&configPath, "c", "", "pipeline YAML file")
	flag.StringVar(&pipeline, "p", "", "pipeline name")
	flag.StringVar(&target, "t", "", "concourse target")
	flag.StringVar(&variablesPath, "l", "", "template values in configuration from a YAML file")
	flag.Parse()

	checkArgument(configPath)
	checkArgument(target)

	action := newReconfigurePipeline()
	err := action.Run(target, pipeline, configPath, variablesPath)

	if err != nil {
		os.Exit(1)
	}
}

func newReconfigurePipeline() *actions.ReconfigurePipeline {
	commandRunner := linux_command_runner.New()

	reconfigurer := concourse.NewReconfigurer(commandRunner)
	processor := lastpass.NewProcessor(commandRunner)
	fifoWriter := fifo.NewWriter()

	return actions.NewReconfigurePipeline(reconfigurer, processor, fifoWriter)
}

func checkArgument(arg string) {
	if arg == "" {
		flag.Usage()
		os.Exit(2)
	}
}
