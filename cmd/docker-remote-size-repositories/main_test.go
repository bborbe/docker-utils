// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Docker Remote Sha for Tag", func() {
	It("Compiles", func() {
		var err error
		_, err = gexec.Build("github.com/bborbe/docker-utils/cmd/docker-remote-size-repositories")
		Expect(err).NotTo(HaveOccurred())
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Docker Remote Sha for Tag Suite")
}
