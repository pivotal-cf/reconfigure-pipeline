package lastpass

import (
	"bytes"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"code.cloudfoundry.org/commandrunner"
	"gopkg.in/yaml.v2"
)

type Processor interface {
	Process(config string) string
}

type processor struct {
	commandRunner commandrunner.CommandRunner
}

func NewProcessor(commandRunner commandrunner.CommandRunner) Processor {
	return &processor{
		commandRunner: commandRunner,
	}
}

func (l *processor) Process(config string) string {
	re := regexp.MustCompile("lpass:///(.*)")

	processedConfig := re.ReplaceAllStringFunc(config, func(match string) string {
		credHandle, _ := url.Parse(match)
		return l.handle(credHandle)
	})

	return processedConfig
}

func (l *processor) handle(credHandle *url.URL) string {
	pathParts := strings.Split(credHandle.Path, "/")

	credential := l.getCredential(pathParts[1], pathParts[2])

	if credHandle.Fragment != "" {
		// Assume YAML contents, return element
		fragmentMap := map[string]string{}
		err := yaml.Unmarshal([]byte(credential), &fragmentMap)
		if err != nil {
			log.Fatal(err)
		}
		credential = fragmentMap[credHandle.Fragment]
	}

	if strings.Contains(credential, "\n") {
		encoded, _ := json.Marshal(credential) // always a string
		return string(encoded)
	}

	return credential
}

func (l *processor) getCredential(credential, field string) string {
	fieldFlagMap := map[string]string{
		"Password": "--password",
		"Username": "--username",
		"URL":      "--url",
		"Notes":    "--notes",
	}

	fieldFlag := fieldFlagMap[field]
	if fieldFlag == "" {
		fieldFlag = "--field=" + field
	}

	output := &bytes.Buffer{}

	cmd := exec.Command("lpass", "show", fieldFlag, credential)

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = output

	err := l.commandRunner.Run(cmd)
	if err != nil {
		log.Fatal(err)
	}

	return output.String()
}
