package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("reconfigure-pipeline", func() {
	It("prints an error message when called with no arguments", func() {
		session := runCommand()
		Eventually(session).Should(Exit(2))

		Expect(session.Err).To(Say("Usage of"))
	})
})
