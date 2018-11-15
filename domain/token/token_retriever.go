package token

import (
	"context"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"encoding/json"
	client "github.com/fabric8-services/fabric8-auth-client/auth"
	"github.com/fabric8-services/build-tool-detector/log"
	"github.com/fabric8-services/fabric8-common/goasupport"
)

// TokenRetriever is used to query auth service and retrieve external providers' token.
type TokenRetriever struct {
	AuthClient *client.Client
	Context    *context.Context
}

// TokenForService calls auth service to retrieve a token for an external service (ie: GitHub).
func (tr *TokenRetriever) TokenForService(forService string) (*string, error) {

	resp, err := tr.AuthClient.RetrieveToken(goasupport.ForwardContextRequestID(*tr.Context), client.RetrieveTokenPath(), forService, nil)
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