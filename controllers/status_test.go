package controllers_test

import (
	"github.com/fabric8-services/build-tool-detector/app/test"
	controllers "github.com/fabric8-services/build-tool-detector/controllers"
	"github.com/goadesign/goa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
)

var _ = Describe("Status", func() {

	Context("Configuration", func() {
		var service *goa.Service

		BeforeEach(func() {
			service = goa.New("build-tool-detector")
		})
		AfterEach(func() {
			gock.Off()
		})

		It("Configuration incorrect - No github_client_id / github_client_secret", func() {
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
