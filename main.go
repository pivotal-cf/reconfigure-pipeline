package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/pivotal-cf/reconfigure-pipeline/actions"
	"github.com/pivotal-cf/reconfigure-pipeline/commandrunner"
	"github.com/pivotal-cf/reconfigure-pipeline/concourse"
	"github.com/pivotal-cf/reconfigure-pipeline/lastpass"
	"github.com/pivotal-cf/reconfigure-pipeline/writer"
)

func main() {
	var opts struct {
		ConfigPath    string `short:"c" long:"config" required:"true" description:"Pipeline configuration file"`
		Pipeline      string `short:"p" long:"pipeline" required:"true" description:"Pipeline to configure"`
		Target        string `short:"t" long:"target" required:"true" description:"Concourse target name"`
		VariablesPath string `short:"l" long:"load-vars-from" description:"Variable flag that can be used for filling in template values in configuration from a YAML file"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(2)
	}

	action := newReconfigurePipeline()
	err = action.Run(opts.Target, opts.Pipeline, opts.ConfigPath, opts.VariablesPath)

	if err != nil {
		os.Exit(1)
	}
}

func newReconfigurePipeline() *actions.ReconfigurePipeline {
	commandRunner := commandrunner.NewSimpleCommandRunner()

	reconfigurer := concourse.NewReconfigurer(commandRunner)
	processor := lastpass.NewProcessor(commandRunner)
	configWriter := writer.NewConfigWriter()

	return actions.NewReconfigurePipeline(reconfigurer, processor, configWriter)
}
