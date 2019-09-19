package lastpass_test

import (
	"fmt"
	"strings"

	"code.cloudfoundry.org/commandrunner/fake_command_runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"

	"errors"

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

	It("identifies credentials", func() {
		commandRunner.WhenRunning(fake_command_runner.CommandSpec{
			Path: "lpass",
			Args: []string{
				"ls",
				"my-folder",
			},
		}, func(cmd *exec.Cmd) error {
			cmd.Stdout.Write([]byte("my-folder\n    my-credential [id: 9999999]"))
			return nil
		})

		input := []string{"my-folder", "my-credential", "Username"}
		expected := lastpass.Credential{
			Name:      "my-folder/my-credential",
			FlagIndex: 2,
		}
		err, output := processor.FindCredential(input)

		Expect(err).ToNot(HaveOccurred())
		Expect(output).To(Equal(expected))
	})

	It("returns error when credential doesn't exist", func() {
		input := []string{"my-folder", "my-credential", "Username"}
		err, _ := processor.FindCredential(input)

		Expect(err).To(Equal(errors.New("credential does not exist")))
	})

	Context("when some folder exist", func() {
		folder := "some-folder"
		cred := "my-credential"

		BeforeEach(func() {
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"ls",
					folder,
				},
			}, func(cmd *exec.Cmd) error {
				cmd.Stdout.Write([]byte(fmt.Sprintf("%s\n    %s [id: 9999999]", folder, cred)))
				return nil
			})
		})

		It("fetches usernames", func() {
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"show",
					"--username",
					fmt.Sprintf("%s/%s", folder, cred),
				},
			}, func(cmd *exec.Cmd) error {
				cmd.Stdout.Write([]byte("my-username"))
				return nil
			})

			input := "key: ((some-folder/my-credential/Username))"
			output := processor.Process(input)

			Expect(output).To(Equal(`key: "my-username"`))
		})

		It("supports fragments for notes", func() {
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"show",
					"--notes",
					fmt.Sprintf("%s/%s", folder, cred),
				},
			}, func(cmd *exec.Cmd) error {
				cmd.Stdout.Write([]byte("inner-key: inner-value\n"))
				return nil
			})

			input := "key: ((some-folder/my-credential/Notes/inner-key))"
			output := processor.Process(input)

			Expect(output).To(Equal(`key: "inner-value"`))
		})
	})

	Context("when credentials exist not in a folder", func() {
		cred := "my-credential"

		BeforeEach(func() {
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"ls",
					cred,
				},
			}, func(cmd *exec.Cmd) error {
				cmd.Stdout.Write([]byte(fmt.Sprintf("((none))\n    %s [id: 9999999]", cred)))
				return nil
			})
		})

		It("fetches usernames", func() {
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"show",
					"--username",
					cred,
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
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"show",
					"--password",
					cred,
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
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"show",
					"--url",
					cred,
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
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"show",
					"--notes",
					cred,
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
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"show",
					"--field=my-field",
					cred,
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
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"show",
					"--notes",
					cred,
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
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"show",
					"--notes",
					cred,
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
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"show",
					"--notes",
					cred,
				},
			}, func(cmd *exec.Cmd) error {
				cmd.Stdout.Write([]byte("inner-key: inner-value\n"))
				return nil
			})

			input := "key: ((my-credential/Notes/inner-key))"
			output := processor.Process(input)

			Expect(output).To(Equal(`key: "inner-value"`))
		})

		It("supports arrays of strings in YAML notes", func() {
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"show",
					"--notes",
					cred,
				},
			}, func(cmd *exec.Cmd) error {
				cmd.Stdout.Write([]byte("inner-key:\n-  inner-value-1\n- inner-value-2\n"))
				return nil
			})

			input := "key: ((my-credential/Notes/inner-key))"
			output := processor.Process(input)

			Expect(output).To(Equal(`key: ["inner-value-1","inner-value-2"]`))
		})

		It("does not call LastPass multiple times for the same credential", func() {
			commandRunner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "lpass",
				Args: []string{
					"show",
					"--notes",
					cred,
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

			callsToShow := getCallCount("show", commandRunner)
			Expect(callsToShow).To(Equal(1))

			Expect(output).To(Equal(`key-1: "inner-value-1"
key-2: "inner-value-2"`))
		})
	})

	It("leaves top level fields alone", func() {
		input := "key: ((top_level_field))"
		output := processor.Process(input)

		Expect(output).To(Equal(`key: ((top_level_field))`))
	})

	It("leaves unknown fields alone", func() {
		commandRunner.WhenRunning(fake_command_runner.CommandSpec{
			Path: "lpass",
			Args: []string{
				"show",
				"--field=secret",
				"unknown",
			},
		}, func(cmd *exec.Cmd) error {
			return errors.New("Exit Status 1")
		})

		input := "key: ((unknown/secret))"
		output := processor.Process(input)

		Expect(output).To(Equal(`key: ((unknown/secret))`))
	})

	It("caches lpass error values", func() {
		commandRunner.WhenRunning(fake_command_runner.CommandSpec{
			Path: "lpass",
			Args: []string{
				"ls",
				"unknown",
			},
		}, func(cmd *exec.Cmd) error {
			cmd.Stdout.Write([]byte(fmt.Sprintf("((none))\n    %s [id: 9999999]", "unknown")))
			return nil
		})

		commandRunner.WhenRunning(fake_command_runner.CommandSpec{
			Path: "lpass",
			Args: []string{
				"show",
				"--field=secret",
				"unknown",
			},
		}, func(cmd *exec.Cmd) error {
			return errors.New("Exit Status 1")
		})

		input := `key-1: ((unknown/secret))
key-2: ((unknown/secret))`
		output := processor.Process(input)

		Expect(output).To(Equal(`key-1: ((unknown/secret))
key-2: ((unknown/secret))`))

		callsToShow := getCallCount("show", commandRunner)
		Expect(callsToShow).To(Equal(1))
	})
})

func getCallCount(command string, runner *fake_command_runner.FakeCommandRunner) int {
	calls := 0
	for _, call := range runner.ExecutedCommands() {
		if strings.Contains(strings.Join(call.Args, ","), command) {
			calls++
		}
	}

	return calls
}
