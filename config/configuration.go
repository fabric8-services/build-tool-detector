/*

Package config implements a way to to
handle configuration management using
viper.

*/
package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	authURI     = "auth.uri"
	serverHost  = "server.host"
	serverPort  = "server.port"
	metricsPort = "server.port"
	sentryDSN   = "sentry.dsn"
)

const (
	defaultAuth = "https://auth.prod-preview.openshift.io"
	defaultHost = "localhost"
	defaultPort = "8099"
)

const (
	prefix       = "BUILD_TOOL_DETECTOR"
	authKeysPath = "/api/token/keys"
)

// Configuration for build tool detector.
type Configuration struct {
	viper *viper.Viper
}

// New returns a configuration with defaults set.
func New() *Configuration {

	// Create new viper
	configuration := Configuration{
		viper: viper.New(),
	}

	// Setup configuration.
	configuration.viper.SetEnvPrefix(prefix)
	configuration.viper.AutomaticEnv()
	configuration.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	configuration.viper.SetTypeByDefaultValue(true)
	configuration.setConfigDefaults()

	return &configuration
}

// GetAuthServiceURL returns the server's port.
func (c *Configuration) GetAuthServiceURL() string {
	return c.viper.GetString(authURI)
}

// GetHost returns the server's host.
func (c *Configuration) GetHost() string {
	return c.viper.GetString(serverHost)
}

// GetPort returns the server's port.
func (c *Configuration) GetPort() string {
	return c.viper.GetString(serverPort)
}

// GetMetricsPort returns the server's port.
func (c *Configuration) GetMetricsPort() string {
	return c.viper.GetString(metricsPort)
}

// GetSentryDSN returs the github client id.
func (c *Configuration) GetSentryDSN() string {
	return c.viper.GetString(sentryDSN)
}

// GetAuthKeysPath provides a URL path to be called for retrieving the keys.
func (c *Configuration) GetAuthKeysPath() string {
	// Fixed with https://github.com/fabric8-services/fabric8-common/pull/25.
	return authKeysPath
}

// GetDevModePrivateKey not used right now.
func (c *Configuration) GetDevModePrivateKey() []byte {
	// No need for now
	return nil
}

// setConfigDefaults sets defaults for configuration.
func (c *Configuration) setConfigDefaults() {
	c.viper.SetDefault(authURI, defaultAuth)
	c.viper.SetDefault(serverHost, defaultHost)
	c.viper.SetDefault(serverPort, defaultPort)
	c.viper.SetDefault(metricsPort, defaultPort)
}
