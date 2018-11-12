/*

Package git_test is used to test the functionality
within the git package.

*/
package repository_test

import (
	"github.com/fabric8-services/build-tool-detector/config"
	"github.com/fabric8-services/build-tool-detector/domain/repository"
	"github.com/fabric8-services/build-tool-detector/domain/repository/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GitServiceType", func() {
	var configuration *config.Configuration

	BeforeSuite(func() {
		configuration = config.New()

	})
	Context("CreateService", func() {
		It("Faulty Host - empty", func() {
			serviceType, err := repository.CreateService("", nil, *configuration)
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrInvalidPath.Error()), "service type should be '400'")
		})

		It("Faulty Host - non-existent", func() {
			serviceType, err := repository.CreateService("test/test", nil, *configuration)
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrInvalidPath.Error()), "service type should be '400'")
		})

		It("Faulty Host - not github.com", func() {
			serviceType, err := repository.CreateService("http://test.com/test/test", nil, *configuration)
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(repository.ErrUnsupportedService.Error()), "service type should be '500'")
		})

		It("Faulty url - no repository", func() {
			serviceType, err := repository.CreateService("http://github.com/test", nil, *configuration)
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrUnsupportedGithubURL.Error()), "service type should be '400'")
		})

		It("Correct url - non-existent", func() {
			serviceType, err := repository.CreateService("http://github.com/test/test", nil, *configuration)
			Expect(serviceType).ShouldNot(BeNil(), "service type should be not be'nil'")
			Expect(err).Should(BeNil(), "err should be 'nil'")
		})
	})

})
