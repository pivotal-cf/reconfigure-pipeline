package actions

import (
	"io/ioutil"
	"log"
	"os"
)

type ReconfigurePipeline struct {
	reconfigurer Reconfigurer
	processor    Processor
	configWriter Writer
}

func NewReconfigurePipeline(
	reconfigurer Reconfigurer,
	processor Processor,
	configWriter Writer,
) *ReconfigurePipeline {
	return &ReconfigurePipeline{
		reconfigurer: reconfigurer,
		processor:    processor,
		configWriter: configWriter,
	}
}

func (r *ReconfigurePipeline) Run(target, pipeline, configPath, variablesPath string) error {
	processedConfigPath, err := r.processConfigFile(configPath)
	defer os.Remove(processedConfigPath)

	if err != nil {
		log.Fatal(err)
	}

	return r.reconfigurer.Reconfigure(target, pipeline, processedConfigPath, variablesPath)
}

func (r *ReconfigurePipeline) processConfigFile(path string) (string, error) {
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	config := string(configBytes)

	processedConfig := r.processor.Process(config)

	return r.configWriter.Write(processedConfig)
}
