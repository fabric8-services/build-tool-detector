/*

Package token_test is used to test the functionality
within the token package.

*/
package token_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTokenRetriever(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GetGitHubToken Suite")
}
