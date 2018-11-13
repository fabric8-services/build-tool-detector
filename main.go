package main

import (
	"net/http"

	"github.com/fabric8-services/build-tool-detector/app"
	"github.com/fabric8-services/build-tool-detector/config"
	"github.com/fabric8-services/build-tool-detector/controllers"
	"github.com/fabric8-services/build-tool-detector/log"
	"github.com/fabric8-services/fabric8-common/goamiddleware"
	"github.com/fabric8-services/fabric8-common/token"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	startup           = "startup"
	errorz            = "err"
	buildToolDetector = "build-tool-detector"
)

func main() {

	// Get a new configuration.
	configuration := config.New()

	service := goa.New(buildToolDetector)

	// Mount middleware.
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	tokenManager, err := token.NewManager(configuration)
	if err != nil {
		log.Logger().Panic(nil, map[string]interface{}{
			"err": err,
		}, "failed to create token manager")
	}
	// Middleware that extracts and stores the token in the context.
	jwtMiddlewareTokenContext := goamiddleware.TokenContext(tokenManager, app.NewJWTSecurity())
	service.Use(jwtMiddlewareTokenContext)

	service.Use(token.InjectTokenManager(tokenManager))
	app.UseJWTMiddleware(service, jwt.New(tokenManager.PublicKeys(), nil, app.NewJWTSecurity()))

	// Mount "build-tool-detector" controller.
	c := controllers.NewBuildToolDetectorController(service, *configuration)
	app.MountBuildToolDetectorController(service, c)

	cs := controllers.NewSwaggerController(service)
	app.MountSwaggerController(service, cs)

	app.MountStatusController(service, controllers.NewStatusController(service))

	// Start/mount metrics http.
	if configuration.GetMetricsPort() == configuration.GetPort() {
		http.Handle("/metrics", promhttp.Handler())
	} else {

		go func(metricAddress string) {
			mx := http.NewServeMux()
			mx.Handle("/metrics", promhttp.Handler())
			if err := http.ListenAndServe(metricAddress, mx); err != nil {
				service.LogError("startup", "err", err)
			}
			// If no port defined in config file, pick one available.
		}(":" + configuration.GetMetricsPort())
	}

	// Start service.
	if err := service.ListenAndServe(":" + configuration.GetPort()); err != nil {
		service.LogError(startup, errorz, err)
	}
}
