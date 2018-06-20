package lastpass

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"code.cloudfoundry.org/commandrunner"
	"gopkg.in/yaml.v2"
	"fmt"
)

type Processor struct {
	commandRunner   commandrunner.CommandRunner
	credentialCache map[string]cacheResult
}

type cacheResult struct {
	 Err error
	 Result string
}

func NewProcessor(commandRunner commandrunner.CommandRunner) *Processor {
	return &Processor{
		commandRunner:   commandRunner,
		credentialCache: map[string]cacheResult{},
	}
}

func (l *Processor) Process(config string) string {
	l.verifyLoggedIn()
	re := regexp.MustCompile(`\(\((.*)\)\)`)

	processedConfig := re.ReplaceAllStringFunc(config, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		return l.handle(submatches[1])
	})

	return processedConfig
}

func (l *Processor) handle(credHandle string) string {
	var encoded []byte
	var err error

	pathParts := strings.Split(credHandle, "/")
	if len(pathParts) == 1 {
	  encoded, _ = json.Marshal(fmt.Sprintf("((%s))", credHandle))
		return fmt.Sprintf("((%s))", credHandle)
	}

	err, credential := l.getCredential(pathParts[0], pathParts[1])
	if err != nil {
		return fmt.Sprintf("((%s))", credHandle)
	}

	fragment := ""
	if len(pathParts) > 2 {
		fragment = pathParts[2]
	}

	if fragment != "" {
		// Assume YAML contents, return element
		fragmentMap := map[string]interface{}{}
		err := yaml.Unmarshal([]byte(credential), &fragmentMap)
		if err != nil {
			log.Fatalln(err)
		}

		value, found := fragmentMap[fragment]
		if !found {
			log.Fatalf("could not find key '%s'\n", fragment)
		}

		encoded, _ = json.Marshal(value)
	} else {
		encoded, _ = json.Marshal(credential)
	}

	return string(encoded)
}

func (l *Processor) getCredential(credential, field string) (error, string) {
	var err error
	cacheKey := strings.Join([]string{credential, field}, "/")
	credentialValue := l.credentialCache[cacheKey].Result
	err = l.credentialCache[cacheKey].Err

	if credentialValue == "" &&  err == nil {
		err, credentialValue = l.getCredentialFromLastPass(credential, field)
		l.credentialCache[cacheKey] = cacheResult{err, credentialValue}
	}

	return err, credentialValue
}

func (l *Processor) getCredentialFromLastPass(credential, field string) (error, string) {
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

	cmd.Stdout = output

	err := l.commandRunner.Run(cmd)
	if err != nil {
		return err, ""
	}

	return nil, strings.TrimSpace(output.String())
}

func (l *Processor) verifyLoggedIn() {
	cmd := exec.Command("lpass", "status")

	output := &bytes.Buffer{}

	cmd.Stdout = output
	cmd.Stderr = output

	err := l.commandRunner.Run(cmd)
	if err != nil {
		log.Fatal(fmt.Sprintf("lpass error: %s", output))
	}
}
