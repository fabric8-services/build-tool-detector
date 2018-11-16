/*

Package git_test is used to test the functionality
within the git package.

*/
package repository_test

import (
	"context"

	"github.com/fabric8-services/build-tool-detector/config"
	"github.com/fabric8-services/build-tool-detector/domain/repository"
	"github.com/fabric8-services/build-tool-detector/domain/repository/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GitServiceType", func() {
	var configuration *config.Configuration
	ctx := context.TODO()

	BeforeSuite(func() {
		configuration = config.New()

	})
	Context("CreateService", func() {
		It("Faulty Host - empty", func() {
			serviceType, err := repository.CreateService(&ctx, "", nil, *configuration)
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrInvalidPath.Error()), "service type should be '400'")
		})

		It("Faulty Host - non-existent", func() {
			serviceType, err := repository.CreateService(&ctx, "test/test", nil, *configuration)
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrInvalidPath.Error()), "service type should be '400'")
		})

		It("Faulty Host - not github.com", func() {
			serviceType, err := repository.CreateService(&ctx, "http://test.com/test/test", nil, *configuration)
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(repository.ErrUnsupportedService.Error()), "service type should be '500'")
		})

		It("Faulty url - no repository", func() {
			serviceType, err := repository.CreateService(&ctx, "http://github.com/test", nil, *configuration)
			Expect(serviceType).Should(BeNil(), "service type should be 'nil'")
			Expect(err.Error()).Should(BeEquivalentTo(github.ErrUnsupportedGithubURL.Error()), "service type should be '400'")
		})

	})

})
