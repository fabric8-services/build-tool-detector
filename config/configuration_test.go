package config_test

import (
	"os"

	"github.com/fabric8-services/build-tool-detector/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/h2non/gock.v1"
)

var _ = Describe("Configuration", func() {
	Context("Configuration Defaults", func() {
		var configuration *config.Configuration

		BeforeEach(func() {
			configuration = config.New()
		})
		AfterEach(func() {
			gock.Off()
		})

		It("Configuration defaults - test defaults are set", func() {
			Expect(configuration.GetHost()).Should(Equal("localhost"), "the host should default to localhost")
			Expect(configuration.GetPort()).Should(Equal("8099"), "the port should default to 8099")
			Expect(configuration.GetMetricsPort()).Should(Equal("8099"), "the metrics port should default to 8099")
			Expect(configuration.GetAuthServiceURL()).Should(Equal("https://auth.prod-preview.openshift.io"), "the auth url should default to https://auth.prod-preview.openshift.io")
			Expect(configuration.GetSentryDSN()).Should(Equal(""), "the sentry dsn should default to empty")
			Expect(configuration.GetAuthKeysPath()).Should(Equal("/api/token/keys"), "the sentry dsn should return /api/token/keys")
		})
	})

	Context("Configuration Overriden Defaults", func() {
		var configuration *config.Configuration

		BeforeEach(func() {
			os.Setenv("BUILD_TOOL_DETECTOR_METRICS_PORT", "1234")
			os.Setenv("BUILD_TOOL_DETECTOR_SERVER_PORT", "1234")
			os.Setenv("BUILD_TOOL_DETECTOR_SERVER_HOST", "test")
			os.Setenv("BUILD_TOOL_DETECTOR_AUTH_URI", "test")
			os.Setenv("BUILD_TOOL_DETECTOR_SENTRY_DSN", "test")
			configuration = config.New()
		})
		AfterEach(func() {
			gock.Off()
			os.Unsetenv("BUILD_TOOL_DETECTOR_METRICS_PORT")
			os.Unsetenv("BUILD_TOOL_DETECTOR_SERVER_PORT")
			os.Unsetenv("BUILD_TOOL_DETECTOR_SERVER_HOST")
			os.Unsetenv("BUILD_TOOL_DETECTOR_AUTH_URI")
			os.Unsetenv("BUILD_TOOL_DETECTOR_SENTRY_DSN")
		})
		It("Configuration defaults - test defaults are overriden", func() {
			Expect(configuration.GetHost()).Should(Equal("test"), "the host should override to test")
			Expect(configuration.GetPort()).Should(Equal("1234"), "the port should override to 1234")
			Expect(configuration.GetMetricsPort()).Should(Equal("1234"), "the metrics port should override to 1234")
			Expect(configuration.GetAuthServiceURL()).Should(Equal("test"), "the auth url should override to test")
			Expect(configuration.GetSentryDSN()).Should(Equal("test"), "the sentry dsn should override to test")
		})
	})
})
