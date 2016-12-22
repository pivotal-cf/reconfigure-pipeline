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
		processor     *lastpass.Processor
		commandRunner *fake_command_runner.FakeCommandRunner
	)

	BeforeEach(func() {
		commandRunner = fake_command_runner.New()

		processor = lastpass.NewProcessor(commandRunner)
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

		input := "key: ((my-credential/Username))"
		output := processor.Process(input)

		Expect(output).To(Equal(`key: "my-username"`))
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

		input := "key: ((my-credential/Password))"
		output := processor.Process(input)

		Expect(output).To(Equal(`key: "my-password"`))
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

		input := "key: ((my-credential/URL))"
		output := processor.Process(input)

		Expect(output).To(Equal(`key: "my-url"`))
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

		input := "key: ((my-credential/Notes))"
		output := processor.Process(input)

		Expect(output).To(Equal(`key: "my-notes"`))
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

		input := "key: ((my-credential/my-field))"
		output := processor.Process(input)

		Expect(output).To(Equal(`key: "my-value"`))
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
			cmd.Stdout.Write([]byte("line-1\nline-2\n"))
			return nil
		})

		input := "key: ((my-credential/Notes))"
		output := processor.Process(input)

		Expect(output).To(Equal(`key: "line-1\nline-2"`))
	})

	It("encodes embedded JSON", func() {
		commandRunner.WhenRunning(CommandSpec{
			Path: "lpass",
			Args: []string{
				"show",
				"--notes",
				"my-credential",
			},
		}, func(cmd *exec.Cmd) error {
			cmd.Stdout.Write([]byte(`{"inner-key":"inner-value"}`))
			return nil
		})

		input := "key: ((my-credential/Notes))"
		output := processor.Process(input)

		Expect(output).To(Equal(`key: "{\"inner-key\":\"inner-value\"}"`))
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

		input := "key: ((my-credential/Notes/inner-key))"
		output := processor.Process(input)

		Expect(output).To(Equal(`key: "inner-value"`))
	})

	It("does not call LastPass multiple times for the same credential", func() {
		commandRunner.WhenRunning(CommandSpec{
			Path: "lpass",
			Args: []string{
				"show",
				"--notes",
				"my-credential",
			},
		}, func(cmd *exec.Cmd) error {
			cmd.Stdout.Write([]byte(`inner-key-1: inner-value-1
inner-key-2: inner-value-2
`))
			return nil
		})

		input := `key-1: ((my-credential/Notes/inner-key-1))
key-2: ((my-credential/Notes/inner-key-2))`

		output := processor.Process(input)

		Expect(commandRunner.ExecutedCommands()).To(HaveLen(1))

		Expect(output).To(Equal(`key-1: "inner-value-1"
key-2: "inner-value-2"`))
	})
})
