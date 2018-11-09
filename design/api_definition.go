/*

Package design is used to develop
the REST endpoints for the build tool.

*/package design

import (
	a "github.com/goadesign/goa/design/apidsl"
)

// API the function to define the top-level API DSL of the application.
var _ = a.API("build-tool-detector", func() {
	a.Title("Build Tool Detector")
	a.Host("openshift.io")
	a.Scheme("http")
	a.BasePath("/api")
	a.Consumes("application/json")
	a.Produces("application/json")

	a.License(func() {
		a.Name("Apache License Version 2.0")
		a.URL("http://www.apache.org/licenses/LICENSE-2.0")
	})
	a.Origin("/[.*openshift.io|localhost]/", func() {
		a.Methods("GET", "POST", "PUT", "PATCH", "DELETE")
		a.Headers("X-Request-Id", "Content-Type", "Authorization", "If-None-Match", "If-Modified-Since")
		a.MaxAge(600)
		a.Credentials()
	})

	a.JWTSecurity("jwt", func() {
		a.Description("JWT Token Auth")
		a.Header("Authorization")
	})

})
