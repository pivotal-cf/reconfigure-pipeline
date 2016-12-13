package actions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"path"

	"github.com/pivotal-cf/reconfigure-pipeline/actions"
	"github.com/pivotal-cf/reconfigure-pipeline/actions/actionsfakes"
)

var _ = Describe("Reconfigure Pipeline", func() {
	var (
		action       *actions.ReconfigurePipeline
		reconfigurer *actionsfakes.FakeReconfigurer
		processor    *actionsfakes.FakeProcessor
		writer       *actionsfakes.FakeWriter

		tempDir    string
		configPath string
	)

	BeforeEach(func() {
		reconfigurer = &actionsfakes.FakeReconfigurer{}
		processor = &actionsfakes.FakeProcessor{}
		writer = &actionsfakes.FakeWriter{}

		action = actions.NewReconfigurePipeline(reconfigurer, processor, writer)

		tempDir, _ = ioutil.TempDir("", "reconfigure-pipeline")
		configPath = path.Join(tempDir, "some-config.yml")
		ioutil.WriteFile(configPath, []byte("My Special Config"), 0600)
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	It("writes the processed configuration file to a pipe", func() {
		processor.ProcessReturns("My Processed Config")

		action.Run("my-special-target", "my-special-pipeline", configPath, "")

		Expect(processor.ProcessCallCount()).To(Equal(1))
		Expect(processor.ProcessArgsForCall(0)).To(Equal("My Special Config"))

		Expect(writer.WriteCallCount()).To(Equal(1))
		Expect(writer.WriteArgsForCall(0)).To(Equal("My Processed Config"))
	})

	It("updates the pipeline using the processed configuration", func() {
		writer.WriteReturns("/tmp/processed.yml", nil)

		action.Run("my-special-target", "my-special-pipeline", configPath, "/tmp/vars.yml")

		Expect(reconfigurer.ReconfigureCallCount()).To(Equal(1))

		target, pipeline, configPath, variablesPath := reconfigurer.ReconfigureArgsForCall(0)
		Expect(target).To(Equal("my-special-target"))
		Expect(pipeline).To(Equal("my-special-pipeline"))
		Expect(configPath).To(Equal("/tmp/processed.yml"))
		Expect(variablesPath).To(Equal("/tmp/vars.yml"))
	})
})
