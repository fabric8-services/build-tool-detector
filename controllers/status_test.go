package controllers_test

import (
	"io/ioutil"

	"github.com/fabric8-services/build-tool-detector/app/test"
	controllers "github.com/fabric8-services/build-tool-detector/controllers"
	"github.com/goadesign/goa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gock.v1"
)

var _ = Describe("Status", func() {

	Context("OK Status", func() {
		var service *goa.Service

		BeforeEach(func() {
			service = goa.New("build-tool-detector")
		})
		AfterEach(func() {
			gock.Off()
		})

		It("Status OK - returns the status commit, build time and start time", func() {
			bodyString, err := ioutil.ReadFile("../controllers/test/mock/localhost/ok_status.json")
			Expect(err).Should(BeNil())

			gock.New("https://test:8099").
				Get("/api/status").
				Reply(200).
				BodyString(string(bodyString))

			test.ShowStatusOK(GinkgoT(), nil, nil, controllers.NewStatusController(service))
		})
	})
})
