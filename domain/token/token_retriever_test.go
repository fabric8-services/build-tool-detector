/*

Package token is used to test the functionality
within the token package.

*/
package token_test

import (
	"context"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gock.v1"

	"github.com/fabric8-services/build-tool-detector/config"
	"github.com/fabric8-services/build-tool-detector/domain/token"
)

var _ = Describe("GetGitHubToken", func() {

	Context("OK Status", func() {
		conf := config.New()
		ctx := context.TODO()

		BeforeEach(func() {
			authBodyString, err := ioutil.ReadFile("../../controllers/test/mock/fabric8_auth_backend/return_token.json")
			Expect(err).Should(BeNil())

			gock.New(conf.GetAuthServiceURL()).
				Get("/api/token").
				Reply(200).
				BodyString(string(authBodyString))
		})
		AfterEach(func() {
			gock.Off()
		})

		It("Status OK - returns the token", func() {
			tr, _ := token.GetGitHubToken(&ctx, *conf)
			Expect(*tr).Should(Equal("ACCESS_TOKEN"), "gh token should match the auth service retirved token")
		})
	})
	Context("Error Status", func() {
		conf := config.New()
		ctx := context.TODO()

		BeforeEach(func() {
			gock.New(conf.GetAuthServiceURL()).
				Get("/api/token").
				Reply(500)
		})
		AfterEach(func() {
			gock.Off()
		})

		It("Auth service internal error", func() {
			tr, err := token.GetGitHubToken(&ctx, *conf)
			Expect(tr).Should(BeNil())
			Expect(err).ShouldNot(BeNil())
		})
	})
})
