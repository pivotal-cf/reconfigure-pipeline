package lastpass

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"code.cloudfoundry.org/commandrunner"
	"gopkg.in/yaml.v2"
)

type Processor struct {
	commandRunner commandrunner.CommandRunner
}

func NewProcessor(commandRunner commandrunner.CommandRunner) *Processor {
	return &Processor{
		commandRunner: commandRunner,
	}
}

func (l *Processor) Process(config string) string {
	re := regexp.MustCompile(`\(\((.*)\)\)`)

	processedConfig := re.ReplaceAllStringFunc(config, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		return l.handle(submatches[1])
	})

	return processedConfig
}

func (l *Processor) handle(credHandle string) string {
	pathParts := strings.Split(credHandle, "/")

	credential := l.getCredential(pathParts[0], pathParts[1])

	fragment := ""
	if len(pathParts) > 2 {
		fragment = pathParts[2]
	}

	if fragment != "" {
		// Assume YAML contents, return element
		fragmentMap := map[string]string{}
		err := yaml.Unmarshal([]byte(credential), &fragmentMap)
		if err != nil {
			log.Fatal(err)
		}
		credential = fragmentMap[fragment]
	}

	if strings.Contains(credential, "\n") {
		encoded, _ := json.Marshal(credential) // always a string
		return string(encoded)
	}

	return credential
}

func (l *Processor) getCredential(credential, field string) string {
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
