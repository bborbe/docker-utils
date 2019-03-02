package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Dockerhub Cleaner", func() {
	It("Compiles", func() {
		var err error
		_, err = gexec.Build("github.com/bborbe/docker-utils/cmd/dockerhub-cleaner")
		Expect(err).NotTo(HaveOccurred())
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dockerhub Cleaner Suite")
}
