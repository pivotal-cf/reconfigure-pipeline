package concourse_test

import (
	. "code.cloudfoundry.org/commandrunner/fake_command_runner"
	. "code.cloudfoundry.org/commandrunner/fake_command_runner/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"os/exec"

	"code.cloudfoundry.org/commandrunner/fake_command_runner"
	"github.com/oozie/reconfigure-pipeline/concourse"
)

var _ = Describe("Reconfigurer", func() {
	var (
		reconfigurer  concourse.Reconfigurer
		commandRunner *fake_command_runner.FakeCommandRunner
	)

	BeforeEach(func() {
		commandRunner = fake_command_runner.New()

		reconfigurer = concourse.NewReconfigurer(commandRunner)
	})

	It("runs fly with the correct arguments", func() {
		err := reconfigurer.Reconfigure("my-target", "my-pipeline", "/tmp/config.yml", "/tmp/vars.yml")
		Expect(err).NotTo(HaveOccurred())

		Expect(commandRunner).To(HaveExecutedSerially(CommandSpec{
			Path: "fly",
			Args: []string{
				"-t", "my-target",
				"set-pipeline",
				"-p", "my-pipeline",
				"-c", "/tmp/config.yml",
				"-l", "/tmp/vars.yml",
			},
		}))
	})

	It("omits -l if variables file is not provided", func() {
		err := reconfigurer.Reconfigure("my-target", "my-pipeline", "/tmp/config.yml", "")
		Expect(err).NotTo(HaveOccurred())

		Expect(commandRunner).To(HaveExecutedSerially(CommandSpec{
			Path: "fly",
			Args: []string{
				"-t", "my-target",
				"set-pipeline",
				"-p", "my-pipeline",
				"-c", "/tmp/config.yml",
			},
		}))
	})

	It("returns an error if fly fails", func() {
		commandRunner.WhenRunning(CommandSpec{
			Path: "fly",
		}, func(cmd *exec.Cmd) error {
			return errors.New("My Special Error")
		})

		err := reconfigurer.Reconfigure("my-target", "my-pipeline", "/tmp/config.yml", "/tmp/vars.yml")
		Expect(err).To(HaveOccurred())
	})
})
