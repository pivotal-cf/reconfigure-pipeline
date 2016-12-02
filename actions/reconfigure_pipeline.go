package actions

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/oozie/reconfigure-pipeline/concourse"
	"github.com/oozie/reconfigure-pipeline/fifo"
	"github.com/oozie/reconfigure-pipeline/lastpass"
)

type ReconfigurePipeline struct {
	reconfigurer *concourse.Reconfigurer
	processor    *lastpass.LastPassProcessor
	fifoWriter   *fifo.Writer
}

func NewReconfigurePipeline(
	reconfigurer *concourse.Reconfigurer,
	processor *lastpass.LastPassProcessor,
	fifoWriter *fifo.Writer,
) *ReconfigurePipeline {
	return &ReconfigurePipeline{
		reconfigurer: reconfigurer,
		processor:    processor,
		fifoWriter:   fifoWriter,
	}
}

func (r *ReconfigurePipeline) Run(target, pipeline, configPath, variablesPath string) error {
	if pipeline == "" {
		pipeline = pipelineNameFromPath(configPath)
	}

	processedConfigPath, err := r.processConfigFile(configPath)
	defer os.Remove(processedConfigPath)

	if err != nil {
		log.Fatal(err)
	}

	return r.reconfigurer.Reconfigure(target, pipeline, configPath, variablesPath)
}

func (r *ReconfigurePipeline) processConfigFile(path string) (string, error) {
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	config := string(configBytes)

	processedConfig := r.processor.Process(config)

	return r.fifoWriter.Write(processedConfig)
}

func pipelineNameFromPath(path string) string {
	foo := filepath.Base(path)

	// Strip the extension
	// TODO: deal with a case of no extension
	parts := strings.Split(foo, ".")

	return strings.Join(parts[0:len(parts)-1], ".")
}
