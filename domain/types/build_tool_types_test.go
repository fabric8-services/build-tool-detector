/*

Package buildtype_test tests
the buildtype package.

*/
package types_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/fabric8-services/build-tool-detector/domain/types"
)

var _ = Describe("BuildToolType", func() {
	Context("Maven", func() {
		It("Get Maven", func() {
			maven := NewMaven()
			Expect(maven.BuildToolType).Should(BeEquivalentTo("maven"), "build tool type type should be 'maven'")
		})
	})

	Context("Unknown", func() {
		It("Get Unknown", func() {
			unknown := NewUnknown()
			Expect(unknown.BuildToolType).Should(BeEquivalentTo("unknown"), "build tool type should be 'unknown'")
		})
	})

	Context("GetTypes", func() {
		It("Get Types", func() {
			nodejs := NewNodeJS()
			maven := NewMaven()
			golang := NewGolang()
			types := GetTypes()

			Expect(types[0].BuildType).Should(BeEquivalentTo(maven.BuildToolType), "build tool type should be 'maven'")
			Expect(types[0].File).Should(BeEquivalentTo("pom.xml"), "file name should be 'pom.xml'")

			Expect(types[1].BuildType).Should(BeEquivalentTo(nodejs.BuildToolType), "build tool type should be 'nodejs'")
			Expect(types[1].File).Should(BeEquivalentTo("package.json"), "file name should be 'package.json'")

			Expect(types[2].BuildType).Should(BeEquivalentTo(golang.BuildToolType), "build tool type should be 'golang'")
			Expect(types[2].File).Should(BeEquivalentTo("main.go"), "file name should be 'main.go'")
		})
	})
})
