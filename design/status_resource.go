/*

Package design is used to develop
the REST endpoints for the build tool.

*/
package design

import (
	d "github.com/goadesign/goa/design"
	a "github.com/goadesign/goa/design/apidsl"
)

var _ = a.Resource("status", func() {
	a.BasePath("/status")
	a.DefaultMedia(Status)
	a.Action("show", func() {
		a.Routing(
			a.GET("/"),
		)
		a.Description("Show the status of the current running instance")
		a.Response(d.OK)
	})
})

// Status defines the status of the current running instance
var Status = a.MediaType("application/vnd.status+json", func() {
	a.Description("The status of the current running instance")
	a.Attributes(func() {
		a.Attribute("commit", d.String, "Commit SHA this build is based on")
		a.Attribute("buildTime", d.String, "The time when built")
		a.Attribute("startTime", d.String, "The time when started")
		a.Attribute("error", d.String, "The error if any")
		a.Required("commit", "buildTime", "startTime")
	})
	a.View("default", func() {
		a.Attribute("commit")
		a.Attribute("buildTime")
		a.Attribute("startTime")
		a.Attribute("error")
	})
})
