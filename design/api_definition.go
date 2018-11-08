/*

Package design is used to develop
the REST endpoints for the build tool.

*/
package design

import (
	"github.com/fabric8-services/build-tool-detector/config"
	"github.com/fabric8-services/build-tool-detector/log"
	a "github.com/goadesign/goa/design/apidsl"
	"github.com/spf13/viper"
	"strconv"
)

var _ = a.API("build-tool-detector", func() {
	a.Title("Build Tool Detector")
	a.Description("Detects the build tool for a specific repository and branch. Currently, this tool only supports detecting the build tool maven for github repositories.")

	a.Origin("/[.*openshift.io|localhost]/", func() {
		a.Methods("GET")
		a.Headers("X-Request-Id", "Content-Type", "Authorization")
		a.MaxAge(600)
		a.Credentials()
	})

	a.Scheme("http")
	a.Host(getHost())
	a.Version("1.0")
	a.BasePath("/build-tool-detector")

	a.License(func() {
		a.Name("Apache License Version 2.0")
		a.URL("http://www.apache.org/licenses/LICENSE-2.0")
	})

	a.JWTSecurity("jwt", func() {
		a.Description("JWT Token Auth")
		a.Header("Authorization")
	})
})

// TODO: Fix this when configuration is fixed: https://github.com/fabric8-services/build-tool-detector/issues/2
func getHost() string {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	var configuration config.Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Logger().Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Logger().Fatalf("unable to decode into struct, %v", err)
	}

	return configuration.Server.Host + strconv.Itoa(configuration.Server.Port)
}
