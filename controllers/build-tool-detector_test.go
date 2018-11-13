/*

Package controllers_test tests the autogenerated
scaffold outputs. Gock is used to mock the
go-github api calls.

*/
package controllers_test

import (
	"io/ioutil"
	"os"

	"github.com/fabric8-services/build-tool-detector/app/test"
	"github.com/fabric8-services/build-tool-detector/config"
	controllers "github.com/fabric8-services/build-tool-detector/controllers"
	"github.com/goadesign/goa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gock.v1"
)

var _ = Describe("BuildToolDetector", func() {

	Context("Configuration", func() {
		var service *goa.Service
		var configuration *config.Configuration

		BeforeEach(func() {
			service = goa.New("build-tool-detector")
			configuration = config.New()
		})
		AfterEach(func() {
			gock.Off()
		})

		It("Configuration incorrect - No github_client_id / github_client_secret", func() {
			bodyString, err := ioutil.ReadFile("../controllers/test/mock/fabric8_launcher_backend/not_found_branch.json")
			Expect(err).Should(BeNil())

			gock.New("https://api.github.com").
				Get("/repos/fabric8-launcher/launcher-backend/branches/master").
				Reply(404).
				BodyString(string(bodyString))

			bodyString, err = ioutil.ReadFile("../controllers/test/mock/fabric8_launcher_backend/not_found_repo_branch.json")
			Expect(err).Should(BeNil())
			gock.New("https://api.github.com").
				Get("/repos/fabric8-launcher/launcher-backend/contents/pom.xml").
				Reply(404).
				BodyString(string(bodyString))
			test.ShowBuildToolDetectorNotFound(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service, *configuration), "https://github.com/fabric8-launcher/launcher-backend/tree/master", nil)
		})
	})

	Context("Internal Server Error", func() {
		var service *goa.Service
		var configuration *config.Configuration

		BeforeEach(func() {
			service = goa.New("build-tool-detector")
			configuration = config.New()
		})
		AfterEach(func() {
			gock.Off()
		})
		It("Non-existent owner name -- 404 Repository Not Found", func() {
			// Fail instead of expect- and make this a method
			bodyString, err := ioutil.ReadFile("../controllers/test/mock/fabric8_launcher_backend/not_found_repo_branch.json")
			Expect(err).Should(BeNil())

			gock.New("https://api.github.com").
				Get("/repos/fabric8-launcherz/launcher-backend/branches/master").
				Reply(404).
				BodyString(string(bodyString))

			branch := "master"
			test.ShowBuildToolDetectorNotFound(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service, *configuration), "https://github.com/fabric8-launcherz/launcher-backend", &branch)
		})

		It("Non-existent owner name -- 404 Owner Not Found", func() {
			bodyString, err := ioutil.ReadFile("../controllers/test/mock/fabric8_launcher_backend/not_found_repo_branch.json")
			Expect(err).Should(BeNil())

			gock.New("https://api.github.com").
				Get("/repos/fabric8-launcher/launcher-backendz/branches/master").
				Reply(404).
				BodyString(string(bodyString))

			branch := "master"
			test.ShowBuildToolDetectorNotFound(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service, *configuration), "https://github.com/fabric8-launcher/launcher-backendz", &branch)
		})

		It("Non-existent branch name -- 404 Branch Not Found", func() {
			bodyString, err := ioutil.ReadFile("../controllers/test/mock/fabric8_launcher_backend/not_found_branch.json")
			Expect(err).Should(BeNil())

			gock.New("https://api.github.com").
				Get("/repos/fabric8-launcher/launcher-backend/branches/masterz").
				Reply(404).
				BodyString(string(bodyString))

			test.ShowBuildToolDetectorNotFound(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service, *configuration), "https://github.com/fabric8-launcher/launcher-backend/tree/masterz", nil)
		})

		It("Invalid URL -- 400 Bad Request", func() {
			branch := "master"
			test.ShowBuildToolDetectorBadRequest(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service, *configuration), "fabric8-launcher/launcher-backend", &branch)
		})

		It("Unsupported Git Service -- 500 Internal Server Error", func() {
			branch := "master"
			test.ShowBuildToolDetectorInternalServerError(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service, *configuration), "http://gitlab.com/fabric8-launcher/launcher-backend", &branch)
		})

		It("Invalid URL and Branch -- 500 Internal Server Error", func() {
			test.ShowBuildToolDetectorBadRequest(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service, *configuration), "", nil)
		})
	})

	Context("Okay", func() {
		var service *goa.Service
		var configuration *config.Configuration

		BeforeEach(func() {
			service = goa.New("build-tool-detector")
			os.Setenv("BUILD_TOOL_DETECTOR_GITHUB_CLIENT_ID", "test")
			os.Setenv("BUILD_TOOL_DETECTOR_GITHUB_CLIENT_SECRET", "test")
			configuration = config.New()
		})
		AfterEach(func() {
			os.Unsetenv("BUILD_TOOL_DETECTOR_GITHUB_CLIENT_ID")
			os.Unsetenv("BUILD_TOOL_DETECTOR_GITHUB_CLIENT_SECRET")
			gock.Off()
		})
		It("Recognize Unknown - Branch field populated", func() {
			bodyString, err := ioutil.ReadFile("../controllers/test/mock/fabric8_wit/ok_branch.json")
			Expect(err).Should(BeNil())

			gock.New("https://api.github.com").
				Get("/repos/fabric8-services/fabric8-wit/branches/master").
				Reply(200).
				BodyString(string(bodyString))

			bodyString, err = ioutil.ReadFile("../controllers/test/mock/fabric8_wit/not_found_contents.json")
			Expect(err).Should(BeNil())
			gock.New("https://api.github.com").
				Get("/repos/fabric8-services/fabric8-wit/contents/pom.xml").
				Reply(404).
				BodyString(string(bodyString))
			branch := "master"
			_, buildTool := test.ShowBuildToolDetectorOK(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service, *configuration), "https://github.com/fabric8-services/fabric8-wit", &branch)
			Expect(buildTool.BuildToolType).Should(Equal("unknown"), "buildTool should not be empty")
		})

		It("Recognize Unknown - Branch included in URL", func() {
			bodyString, err := ioutil.ReadFile("../controllers/test/mock/fabric8_wit/ok_branch.json")
			Expect(err).Should(BeNil())

			gock.New("https://api.github.com").
				Get("/repos/fabric8-services/fabric8-wit/branches/master").
				Reply(200).
				BodyString(string(bodyString))

			bodyString, err = ioutil.ReadFile("../controllers/test/mock/fabric8_wit/not_found_contents.json")
			Expect(err).Should(BeNil())
			gock.New("https://api.github.com").
				Get("/repos/fabric8-services/fabric8-wit/contents/pom.xml").
				Reply(404).
				BodyString(string(bodyString))
			_, buildTool := test.ShowBuildToolDetectorOK(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service, *configuration), "https://github.com/fabric8-services/fabric8-wit/tree/master", nil)
			Expect(buildTool.BuildToolType).Should(Equal("unknown"), "buildTool should not be empty")
		})

		It("Recognize Maven - Branch field populated", func() {
			bodyString, err := ioutil.ReadFile("../controllers/test/mock/fabric8_launcher_backend/ok_branch.json")
			Expect(err).Should(BeNil())

			gock.New("https://api.github.com").
				Get("/repos/fabric8-launcher/launcher-backend/branches/master").
				Reply(200).
				BodyString(string(bodyString))

			bodyString, err = ioutil.ReadFile("../controllers/test/mock/fabric8_launcher_backend/ok_contents.json")
			Expect(err).Should(BeNil())
			gock.New("https://api.github.com").
				Get("/repos/fabric8-launcher/launcher-backend/contents/pom.xml").
				Reply(200).
				BodyString(string(bodyString))
			branch := "master"
			_, buildTool := test.ShowBuildToolDetectorOK(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service, *configuration), "https://github.com/fabric8-launcher/launcher-backend", &branch)
			Expect(buildTool.BuildToolType).Should(Equal("maven"), "buildTool should not be empty")
		})

		It("Recognize Maven - Branch included in URL", func() {
			bodyString, err := ioutil.ReadFile("../controllers/test/mock/fabric8_launcher_backend/ok_branch.json")
			Expect(err).Should(BeNil())

			gock.New("https://api.github.com").
				Get("/repos/fabric8-launcher/launcher-backend/branches/master").
				Reply(200).
				BodyString(string(bodyString))

			bodyString, err = ioutil.ReadFile("../controllers/test/mock/fabric8_launcher_backend/ok_contents.json")
			Expect(err).Should(BeNil())
			gock.New("https://api.github.com").
				Get("/repos/fabric8-launcher/launcher-backend/contents/pom.xml").
				Reply(200).
				BodyString(string(bodyString))
			_, buildTool := test.ShowBuildToolDetectorOK(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service, *configuration), "https://github.com/fabric8-launcher/launcher-backend/tree/master", nil)
			Expect(buildTool.BuildToolType).Should(Equal("maven"), "buildTool should not be empty")
		})

		It("Recognize NodeJS - Branch included in URL", func() {
			bodyString, err := ioutil.ReadFile("../controllers/test/mock/fabric8_ui/ok_branch.json")
			Expect(err).Should(BeNil())

			gock.New("https://api.github.com").
				Get("/repos/fabric8-ui/fabric8-ui/branches/master").
				Reply(200).
				BodyString(string(bodyString))

			bodyString, err = ioutil.ReadFile("../controllers/test/mock/fabric8_ui/ok_contents.json")
			Expect(err).Should(BeNil())
			gock.New("https://api.github.com").
				Get("/repos/fabric8-ui/fabric8-ui/contents/package.json").
				Reply(200).
				BodyString(string(bodyString))
			_, buildTool := test.ShowBuildToolDetectorOK(GinkgoT(), nil, nil, controllers.NewBuildToolDetectorController(service, *configuration), "https://github.com/fabric8-ui/fabric8-ui/tree/master", nil)
			Expect(buildTool.BuildToolType).Should(Equal("nodejs"), "buildTool should be nodejs")
		})
	})
})
