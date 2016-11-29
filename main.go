package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/oozie/reconfigure-pipeline/lastpass"
)

func main() {
	// options
	var pipelineName, configPath, target, variablesFile string
	flag.StringVar(&configPath, "c", "", "pipeline YAML file")
	flag.StringVar(&pipelineName, "p", "", "pipeline name")
	flag.StringVar(&target, "t", "", "concourse target")
	flag.StringVar(&variablesFile, "l", "", "template values in configuration from a YAML file")
	flag.Parse()

	checkArgument(configPath)
	checkArgument(target)

	if pipelineName == "" {
		pipelineName = PipelineNameFromPath(configPath)
	}

	tempConfigPath, err := WriteConfigFile(configPath)
	defer os.Remove(tempConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	args := []string{"-t", target, "set-pipeline", "-p", pipelineName, "-c", tempConfigPath}
	if variablesFile != "" {
		args = append(args, "-l", variablesFile)
	}
	cmd := exec.Command("fly", args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		exitError := err.(*exec.ExitError)
		status, _ := exitError.Sys().(syscall.WaitStatus)
		os.Exit(status.ExitStatus())
	}
}

func PipelineNameFromPath(path string) string {
	foo := filepath.Base(path)

	// Strip the extension
	// TODO: deal with a case of no extension
	parts := strings.Split(foo, ".")

	return strings.Join(parts[0:len(parts)-1], ".")
}

func WriteConfigFile(path string) (string, error) {
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	config := string(configBytes)

	re := regexp.MustCompile("lpass:///(.*)")
	processedConfig := re.ReplaceAllStringFunc(config, func(match string) string {
		credHandle, _ := url.Parse(match)
		return lastpass.Handle(credHandle)
	})

	tempDir, err := ioutil.TempDir("", "reconfigure-pipeline")
	if err != nil {
		return "", err
	}

	fifoPath := filepath.Join(tempDir, "fifo")
	err = syscall.Mkfifo(fifoPath, 0600)
	if err != nil {
		return "", err
	}

	go func() {
		f, err := os.OpenFile(fifoPath, os.O_WRONLY, 0600)
		defer f.Close()

		if err != nil {
			log.Fatal(err)
		}

		f.WriteString(processedConfig)
	}()

	return fifoPath, nil
}

func checkArgument(arg string) {
	if arg == "" {
		flag.Usage()
		os.Exit(2)
	}
}
