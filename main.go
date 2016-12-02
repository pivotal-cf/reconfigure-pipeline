package main

import (
	"flag"
	"os"

	"github.com/oozie/reconfigure-pipeline/actions"
	"github.com/oozie/reconfigure-pipeline/concourse"
	"github.com/oozie/reconfigure-pipeline/fifo"
	"github.com/oozie/reconfigure-pipeline/lastpass"
)

func main() {
	// options
	var pipeline, configPath, target, variablesPath string

	flag.StringVar(&configPath, "c", "", "pipeline YAML file")
	flag.StringVar(&pipeline, "p", "", "pipeline name")
	flag.StringVar(&target, "t", "", "concourse target")
	flag.StringVar(&variablesPath, "l", "", "template values in configuration from a YAML file")
	flag.Parse()

	checkArgument(configPath)
	checkArgument(target)

	reconfigurer := concourse.NewReconfigurer()
	processor := lastpass.NewProcessor()
	fifoWriter := fifo.NewWriter()

	action := actions.NewReconfigurePipeline(reconfigurer, processor, fifoWriter)
	err := action.Run(target, pipeline, configPath, variablesPath)

	if err != nil {
		os.Exit(1)
	}
}

func checkArgument(arg string) {
	if arg == "" {
		flag.Usage()
		os.Exit(2)
	}
}
