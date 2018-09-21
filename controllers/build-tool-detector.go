/*

Package controllers is autogenerated
and containing scaffold outputs
as well as manually created sub-packages
and files.

*/
package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/tinakurian/build-tool-detector/controllers/git"
	"net/http"

	"github.com/goadesign/goa"
	"github.com/tinakurian/build-tool-detector/app"
	"github.com/tinakurian/build-tool-detector/controllers/buildtype"
	errs "github.com/tinakurian/build-tool-detector/controllers/error"
	"github.com/tinakurian/build-tool-detector/controllers/service"
)

// BuildToolDetectorController implements the build-tool-detector resource.
type BuildToolDetectorController struct {
	*goa.Controller
}

// NewBuildToolDetectorController creates a build-tool-detector controller.
func NewBuildToolDetectorController(service *goa.Service) *BuildToolDetectorController {
	return &BuildToolDetectorController{Controller: service.NewController("BuildToolDetectorController")}
}

// Show runs the show action.
func (c *BuildToolDetectorController) Show(ctx *app.ShowBuildToolDetectorContext) error {

	gitService := service.System{}.GetGitService()

	_, err := git.GetGitServiceType(ctx.URL)
	if err != nil {
		return handleRequest(ctx, err, nil)
	}

	err, buildTool := gitService.GetGitHubService().GetContents(ctx)
	if err != nil {
		if err.StatusCode == http.StatusBadRequest {
			return handleRequest(ctx, err, nil)
		}
		return handleRequest(ctx, err, buildtype.Unknown())
	}

	return handleRequest(ctx, nil, buildTool)
}

func handleRequest(ctx *app.ShowBuildToolDetectorContext, httpTypeError *errs.HTTPTypeError, buildTool *app.GoaBuildToolDetector) error {
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")

	if httpTypeError == nil || httpTypeError.StatusCode == http.StatusInternalServerError {
		if buildTool != nil {
			return ctx.OK(buildTool)
		}
	}

	ctx.WriteHeader(httpTypeError.StatusCode)
	jsonHTTPTypeError, err := json.Marshal(httpTypeError)
	if err != nil {
		// TODO: log and return error
		panic(err)
	}

	if _, err := fmt.Fprint(ctx.ResponseWriter, string(jsonHTTPTypeError)); err != nil {
		// TODO: log and return error
		panic(err)
	}

	return getErrResponse(ctx, httpTypeError)
}

func getErrResponse(ctx *app.ShowBuildToolDetectorContext, httpTypeError *errs.HTTPTypeError) error {
	var response error
	switch httpTypeError.StatusCode {
	case http.StatusBadRequest:
		response = ctx.BadRequest()
	case http.StatusInternalServerError:
		response = ctx.InternalServerError()
	}

	return response
}
