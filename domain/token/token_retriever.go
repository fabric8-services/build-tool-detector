package token

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fabric8-services/build-tool-detector/config"
	"github.com/fabric8-services/build-tool-detector/log"
	client "github.com/fabric8-services/fabric8-auth-client/auth"
	"github.com/fabric8-services/fabric8-common/goasupport"
	goaclient "github.com/goadesign/goa/client"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const githubHost = "github.com"

// TokenForService calls auth service to retrieve a token for an external service (ie: GitHub).
func tokenForService(ctx *context.Context, authClient *client.Client, forService string) (*string, error) {

	resp, err := authClient.RetrieveToken(goasupport.ForwardContextRequestID(*ctx), client.RetrieveTokenPath(), forService, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve token")
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	status := resp.StatusCode
	if status != http.StatusOK {
		log.Logger().Error(nil, map[string]interface{}{
			"err":          err,
			"request_path": client.ShowUserPath(),
			"for_service":  forService,
			"http_status":  status,
		}, "failed to GET token from auth service due to HTTP error %s", status)
		return nil, errors.Wrap(err, "failed to GET token from auth service due to HTTP error")
	}

	var respType client.TokenData
	err = json.Unmarshal(respBody, &respType)
	if err != nil {
		log.Logger().Error(nil, map[string]interface{}{
			"err":           err,
			"request_path":  client.ShowUserPath(),
			"for_service":   forService,
			"http_status":   status,
			"response_body": respBody,
		}, "unable to unmarshal Auth token")
		return nil, errors.Wrap(err, "unable to unmarshal Auth token")
	}

	return respType.AccessToken, nil
}

// GetGitHubToken retrieve GitHub token associated to given openshift.io token using auth service.
func GetGitHubToken(ctx *context.Context, configuration config.Configuration) (*string, error) {
	url, err := url.Parse(configuration.GetAuthServiceURL())
	if err != nil {
		return nil, errors.Wrap(err, "auth service url not found")
	}
	authClient := client.New(goaclient.HTTPClientDoer(http.DefaultClient))
	authClient.Host = url.Host
	authClient.Scheme = url.Scheme
	if goajwt.ContextJWT(*ctx) != nil {
		authClient.SetJWTSigner(goasupport.NewForwardSigner(*ctx))
	} else {
		log.Logger().Info(ctx, nil, "no token in context")
	}

	forService := fmt.Sprintf("%s%s://%s", authClient.Host, authClient.Scheme, githubHost)
	token, err := tokenForService(ctx, authClient, forService)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve token from auth")
	}
	return token, nil
}