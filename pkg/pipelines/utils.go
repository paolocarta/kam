package pipelines

import (
	"fmt"
	"net/url"

	"github.com/jenkins-x/go-scm/scm/factory"
)

var defaultClientFactory = factory.FromRepoURL

// BootstrapRepository creates an authenticated request.
func BootstrapRepository(o *BootstrapOptions) error {
	return nil
}

func repoURL(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("failed to parse %q: %w", u, err)
	}
	parsed.Path = ""
	parsed.User = nil
	return parsed.String(), nil
}
