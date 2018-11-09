package controllers

import (
	"github.com/fabric8-services/build-tool-detector/app"
	"github.com/goadesign/goa"
)

// StatusController implements the status resource.
type StatusController struct {
	*goa.Controller
}

// NewStatusController creates a status controller.
func NewStatusController(service *goa.Service) *StatusController {
	return &StatusController{Controller: service.NewController("StatusController")}
}

// Show runs the show action.
func (c *StatusController) Show(ctx *app.ShowStatusContext) error {
	return ctx.OK(&app.Status{Commit: app.Commit, BuildTime: app.BuildTime, StartTime: app.StartTime})
}
