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
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/fabric8-services/build-tool-detector/config"
	"github.com/fabric8-services/build-tool-detector/domain/repository/github"
	"github.com/fabric8-services/build-tool-detector/domain/token"
	"github.com/fabric8-services/build-tool-detector/domain/types"
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

	tk, err := token.GetGitHubToken(ctx, configuration.GetAuthServiceURL(), u)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve token from auth")
	}
	if tk == nil {
		return nil, errors.Wrap(err, "failed to retrieve token from auth")
	}
	return github.Create(urlSegments, branch, configuration, *tk)
}
