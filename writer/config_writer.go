package writer

import (
	"io/ioutil"
	"log"
	"os"
)

type ConfigWriter struct {
	filePath string
}

func NewConfigWriter() *ConfigWriter {
	return &ConfigWriter{}
}

func (cf *ConfigWriter) Write(content string) (string, error) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		return "", err
	}

	configFile, err := ioutil.TempFile(tmpDir, "reconfigure-pipeline")
	if err != nil {
		return "", err
	}

	configFile.Write([]byte(content))
	defer configFile.Close()

	if err != nil {
		log.Fatal(err)
	}
	cf.filePath = configFile.Name()
	return configFile.Name(), nil
}
