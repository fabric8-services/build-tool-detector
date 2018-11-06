/*

Package git_test is used to test the functionality
within the error package.

*/
package error_test

import (
	"errors"
	"net/http"

	. "github.com/fabric8-services/build-tool-detector/controllers/error"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error", func() {
	Context("ErrBadRequest", func() {
		It("Set ErrBadRequest", func() {
			badRequest := ErrBadRequest(errors.New("bad request"))
			Expect(badRequest.StatusCode).Should(BeEquivalentTo(http.StatusBadRequest), "service type should be 'nil'")
		})
	})

	Context("ErrInternalServerError", func() {
		It("Set ErrInternalServerError", func() {
			internalservererror := ErrInternalServerError(errors.New("internal server error"))
			Expect(internalservererror.StatusCode).Should(BeEquivalentTo(http.StatusInternalServerError), "service type should be 'nil'")
		})
	})

	Context("ErrNotFoundError", func() {
		It("Set ErrNotFoundError", func() {
			notfound := ErrNotFoundError(errors.New("not found"))
			Expect(notfound.StatusCode).Should(BeEquivalentTo(http.StatusNotFound), "service type should be 'nil'")
		})
	})
})
