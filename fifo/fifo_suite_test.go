package fifo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFifo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fifo Suite")
}
