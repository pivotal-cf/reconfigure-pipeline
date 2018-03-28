package writer_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"

	"github.com/pivotal-cf/reconfigure-pipeline/writer"
)

var _ = Describe("Config Writer", func() {
	var (
		configWriter *writer.ConfigWriter
	)

	BeforeEach(func() {
		configWriter = writer.NewConfigWriter()
	})

	It("writes to a tmp file and returns the path", func() {
		contents := "My Special Contents"

		path, err := configWriter.Write(contents)
		Expect(err).NotTo(HaveOccurred())

		defer os.Remove(path)

		readContents, err := ioutil.ReadFile(path)
		Expect(err).NotTo(HaveOccurred())

		Expect(readContents).To(Equal([]byte(contents)))
	})
})
