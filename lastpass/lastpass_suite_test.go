package lastpass_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLastpass(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lastpass Suite")
}
