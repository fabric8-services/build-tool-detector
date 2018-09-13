package controllers

import (
	"build-tool-detector/app"
	"build-tool-detector/controllers/git"
	"build-tool-detector/controllers/git/buildtype"
	"build-tool-detector/controllers/git/github"
	"encoding/json"
	"fmt"
	"github.com/goadesign/goa"
	"net/http"
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
	gitType, url := git.GetType(ctx.URL)
	switch gitType {
	case git.GITHUB:
		return handleGitHub(ctx, url)
	case git.UNKNOWN:
		var buildTool app.GoaBuildToolDetector
		return handleUnknown(ctx, buildTool)
	default:
		var buildTool app.GoaBuildToolDetector
		return handleUnknown(ctx, buildTool)
	}
}

func handleGitHub(ctx *app.ShowBuildToolDetectorContext, url []string) error {
	statusCode, buildTool := github.DetectBuildTool(ctx, url)
	if statusCode == http.StatusInternalServerError {
		return handleUnknown(ctx, buildTool)
	}

	return ctx.OK(&buildTool)
}

func handleUnknown(ctx *app.ShowBuildToolDetectorContext, buildTool app.GoaBuildToolDetector) error {
	if buildTool.BuildToolType == "" {
		buildTool = buildtype.Unknown()
	}

	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	ctx.WriteHeader(http.StatusInternalServerError)
	jsonBuildTool, err := json.Marshal(buildTool)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(ctx.ResponseWriter, string(jsonBuildTool))
	return ctx.InternalServerError()
}
