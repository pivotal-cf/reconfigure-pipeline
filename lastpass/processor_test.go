package lastpass_test

import (
	. "code.cloudfoundry.org/commandrunner/fake_command_runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"

	"code.cloudfoundry.org/commandrunner/fake_command_runner"
	"github.com/pivotal-cf/reconfigure-pipeline/lastpass"
)

var _ = Describe("Processor", func() {
	var (
		processor     lastpass.Processor
		commandRunner *fake_command_runner.FakeCommandRunner
	)

	BeforeEach(func() {
		commandRunner = fake_command_runner.New()

		processor = lastpass.NewProcessor(commandRunner)
	})

	It("does not modify a pipeline without LastPass credentials", func() {
		input := "key: credhub:///my-credential"
		output := processor.Process(input)

		Expect(output).To(Equal(output))
	})

	It("fetches usernames", func() {
		commandRunner.WhenRunning(CommandSpec{
			Path: "lpass",
			Args: []string{
				"show",
				"--username",
				"my-credential",
			},
		}, func(cmd *exec.Cmd) error {
			cmd.Stdout.Write([]byte("my-username"))
			return nil
		})

		input := "key: lpass:///my-credential/Username"
		output := processor.Process(input)

		Expect(output).To(Equal("key: my-username"))
	})

	It("fetches passwords", func() {
		commandRunner.WhenRunning(CommandSpec{
			Path: "lpass",
			Args: []string{
				"show",
				"--password",
				"my-credential",
			},
		}, func(cmd *exec.Cmd) error {
			cmd.Stdout.Write([]byte("my-password"))
			return nil
		})

		input := "key: lpass:///my-credential/Password"
		output := processor.Process(input)

		Expect(output).To(Equal("key: my-password"))
	})

	It("fetches URLs", func() {
		commandRunner.WhenRunning(CommandSpec{
			Path: "lpass",
			Args: []string{
				"show",
				"--url",
				"my-credential",
			},
		}, func(cmd *exec.Cmd) error {
			cmd.Stdout.Write([]byte("my-url"))
			return nil
		})

		input := "key: lpass:///my-credential/URL"
		output := processor.Process(input)

		Expect(output).To(Equal("key: my-url"))
	})

	It("fetches notes", func() {
		commandRunner.WhenRunning(CommandSpec{
			Path: "lpass",
			Args: []string{
				"show",
				"--notes",
				"my-credential",
			},
		}, func(cmd *exec.Cmd) error {
			cmd.Stdout.Write([]byte("my-notes"))
			return nil
		})

		input := "key: lpass:///my-credential/Notes"
		output := processor.Process(input)

		Expect(output).To(Equal("key: my-notes"))
	})

	It("fetches other fields", func() {
		commandRunner.WhenRunning(CommandSpec{
			Path: "lpass",
			Args: []string{
				"show",
				"--field=my-field",
				"my-credential",
			},
		}, func(cmd *exec.Cmd) error {
			cmd.Stdout.Write([]byte("my-value"))
			return nil
		})

		input := "key: lpass:///my-credential/my-field"
		output := processor.Process(input)

		Expect(output).To(Equal("key: my-value"))
	})

	It("encodes multi-line strings", func() {
		commandRunner.WhenRunning(CommandSpec{
			Path: "lpass",
			Args: []string{
				"show",
				"--notes",
				"my-credential",
			},
		}, func(cmd *exec.Cmd) error {
			cmd.Stdout.Write([]byte("line-1\nline-2"))
			return nil
		})

		input := "key: lpass:///my-credential/Notes"
		output := processor.Process(input)

		Expect(output).To(Equal(`key: "line-1\nline-2"`))
	})

	It("supports fragments for notes", func() {
		commandRunner.WhenRunning(CommandSpec{
			Path: "lpass",
			Args: []string{
				"show",
				"--notes",
				"my-credential",
			},
		}, func(cmd *exec.Cmd) error {
			cmd.Stdout.Write([]byte("inner-key: inner-value\n"))
			return nil
		})

		input := "key: lpass:///my-credential/Notes#inner-key"
		output := processor.Process(input)

		Expect(output).To(Equal("key: inner-value"))
	})
})
