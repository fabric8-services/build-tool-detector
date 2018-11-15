/*

Package repository handles detecting build tool types
for git services such as github, bitbucket
and gitlab.

Currently the build-tool-detector only
supports github and can only recognize
maven.

*/
package repository

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/fabric8-services/build-tool-detector/config"
	tr "github.com/fabric8-services/build-tool-detector/domain/token"
	"github.com/fabric8-services/build-tool-detector/domain/repository/github"
	"github.com/fabric8-services/build-tool-detector/domain/types"
	client "github.com/fabric8-services/fabric8-auth-client/auth"
	"github.com/fabric8-services/fabric8-common/goasupport"
	goaclient "github.com/goadesign/goa/client"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"
	"github.com/fabric8-services/build-tool-detector/log"
)

var (
	// ErrUnsupportedService git service unsupported.
	ErrUnsupportedService = errors.New("unsupported service")
)

const (
	slash      = "/"
	githubHost = "github.com"
)

// CreateService performs a simple url parse and split
// in order to retrieve the owner, repository
// and potentially the branch.
//
// Note: This method will likely need to be enhanced
// to handle different github url formats.
func CreateService(ctx *context.Context, urlToParse string, branch *string, configuration config.Configuration) (types.RepositoryService, error) {

	u, err := url.Parse(urlToParse)

	// Fail on error or empty host or empty scheme.
	if err != nil || u.Host == "" || u.Scheme == "" {
		return nil, github.ErrInvalidPath
	}

	// Currently only support Github.
	if u.Host != githubHost {
		return nil, ErrUnsupportedService
	}

	urlSegments := strings.Split(u.Path, slash)
	if len(urlSegments) < 3 {
		return nil, github.ErrUnsupportedGithubURL
	}

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
	tokenRetriever := tr.TokenRetriever{AuthClient: authClient, Context: ctx}
	token, err := tokenRetriever.TokenForService(fmt.Sprintf("https://%s", githubHost))
	if err != nil {
		return nil, errors.Wrap(err, "auth service url not found")
	}
	if token == nil {
		return nil, errors.Wrap(err, "token not found for GitHub")
	}
	return github.Create(urlSegments, branch, configuration, *token)
}
