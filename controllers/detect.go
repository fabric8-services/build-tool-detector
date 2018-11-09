/*

Package controllers is autogenerated
and containing scaffold outputs
as well as manually created sub-packages
and files.

*/
package controllers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fabric8-services/build-tool-detector/app"
	"github.com/fabric8-services/build-tool-detector/config"
	errs "github.com/fabric8-services/build-tool-detector/controllers/error"
	"github.com/fabric8-services/build-tool-detector/domain/repository"
	"github.com/fabric8-services/build-tool-detector/domain/repository/github"
	"github.com/fabric8-services/build-tool-detector/domain/types"
	"github.com/fabric8-services/build-tool-detector/log"
	"github.com/goadesign/goa"
)

var (
	// ErrFailedJSONMarshal unable to marshal json.
	ErrFailedJSONMarshal = errors.New("unable to marshal json")

	// ErrFailedPropagate unable to propagate error.
	ErrFailedPropagate = errors.New("unable to propagate error")
)

const (
	errorz                      = "error"
	contentType                 = "Content-Type"
	applicationJSON             = "application/json"
	buildToolDetectorController = "BuildToolDetectorController"
)

// DetectController implements the detect resource.
type DetectController struct {
	*goa.Controller
	config.Configuration
}

// NewDetectController creates a detect controller.
func NewDetectController(service *goa.Service, configuration config.Configuration) *DetectController {
	return &DetectController{Controller: service.NewController(buildToolDetectorController), Configuration: configuration}
}

// Show runs the show action.
func (c *DetectController) Show(ctx *app.ShowDetectContext) error {
	rawURL := ctx.URL
	repositoryService, err := repository.CreateService(rawURL, ctx.Branch, c.Configuration)
	ctx.ResponseWriter.Header().Set(contentType, applicationJSON)
	if err != nil {
		return handleError(ctx, err)
	}

	buildToolType, err := repositoryService.DetectBuildTool(ctx.Context)
	if err != nil {
		return handleError(ctx, err)
	}

	buildTool := handleSuccess(*buildToolType)
	return ctx.OK(buildTool)
}

// handleSuccess handles returning
// the correct json for 200 OK responses.
func handleSuccess(buildToolType string) *app.GoaDetect {
	switch buildToolType {
	case types.Maven:
		return types.NewMaven()
	case types.NodeJS:
		return types.NewNodeJS()
	case types.Unknown:
		return types.NewUnknown()
	default:
		return types.NewUnknown()
	}
}

// handleError handles returning
// the correct http responses upon error.
func handleError(ctx *app.ShowDetectContext, err error) error {
	switch err.Error() {
	case github.ErrInvalidPath.Error():
		httpError := errs.ErrBadRequest(err)
		writerErr := formatResponse(ctx, httpError)
		if writerErr != nil {
			return writerErr
		}
		return ctx.BadRequest()
	case github.ErrResourceNotFound.Error():
		httpError := errs.ErrNotFoundError(err)
		writerErr := formatResponse(ctx, httpError)
		if writerErr != nil {
			return writerErr
		}
		return ctx.NotFound()
	case repository.ErrUnsupportedService.Error(),
		github.ErrUnsupportedGithubURL.Error():
		httpError := errs.ErrInternalServerError(err)
		writerErr := formatResponse(ctx, httpError)
		if writerErr != nil {
			return writerErr
		}
		return ctx.InternalServerError()
	case github.ErrFailedContentRetrieval.Error():
		buildTool := types.NewUnknown()
		return ctx.OK(buildTool)
	default:
		return ctx.InternalServerError()
	}
}

// formatResponse writes the header
// and formats the response.
func formatResponse(ctx *app.ShowDetectContext, httpTypeError *errs.HTTPTypeError) error {
	ctx.WriteHeader(httpTypeError.StatusCode)
	jsonHTTPTypeError, err := json.Marshal(httpTypeError)
	if err != nil {
		log.Logger().WithError(err).WithField(errorz, httpTypeError).Errorf(ErrFailedJSONMarshal.Error())
		return ctx.InternalServerError()
	}

	if _, err := fmt.Fprint(ctx.ResponseWriter, string(jsonHTTPTypeError)); err != nil {
		log.Logger().WithError(err).WithField(errorz, jsonHTTPTypeError).Errorf(ErrFailedPropagate.Error())
		return ctx.InternalServerError()
	}
	return nil
}
