/*

Package token is used to test the functionality
within the token package.

*/
package token_test

import (
	"context"
	"io/ioutil"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gock.v1"

	"github.com/fabric8-services/build-tool-detector/domain/token"
)

var _ = Describe("GetGitHubToken", func() {

	Context("OK Status", func() {
		authURL := "https://auth.prod-preview.openshift.io"
		ctx := context.TODO()
		u, _ := url.Parse("https://github.com")

		BeforeEach(func() {
			authBodyString, err := ioutil.ReadFile("../../controllers/test/mock/fabric8_auth_backend/return_token.json")
			Expect(err).Should(BeNil())

			gock.New(authURL).
				Get("/api/token").
				Reply(200).
				BodyString(string(authBodyString))
		})
		AfterEach(func() {
			gock.Off()
		})

		It("Status OK - returns the token", func() {
			tr, _ := token.GetGitHubToken(&ctx, authURL, u)
			Expect(*tr).Should(Equal("ACCESS_TOKEN"), "gh token should match the auth service retirved token")
		})
	})
	Context("Error Status", func() {
		authURL := "https://auth.prod-preview.openshift.io"
		ctx := context.TODO()

		BeforeEach(func() {
			gock.New(authURL).
				Get("/api/token").
				Reply(500)
		})
		AfterEach(func() {
			gock.Off()
		})

		It("Status 500 - Auth service internal error", func() {
			u, _ := url.Parse("https://github.com")
			tr, err := token.GetGitHubToken(&ctx, authURL, u)
			Expect(tr).Should(BeNil())
			Expect(err).ShouldNot(BeNil())
			Expect(err.Error()).Should(Equal("failed to retrieve token from auth: failed to GET token from auth service due to HTTP error"))
		})
	})
})
