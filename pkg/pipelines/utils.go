package pipelines

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
)

var defaultClientFactory = factory.FromRepoURL

const defaultRepoDescription = "Bootstrapped GitOps Repository"

// BootstrapRepository creates an authenticated request.
func BootstrapRepository(o *BootstrapOptions) error {
	u, err := url.Parse(o.GitOpsRepoURL)
	if err != nil {
		return fmt.Errorf("failed to parse GitOps repo URL %q: %w", o.GitOpsRepoURL, err)
	}
	parts := strings.Split(u.Path, "/")
	ri := &scm.RepositoryInput{
		Private:     true,
		Description: defaultRepoDescription,
		Namespace:   parts[1],
		Name:        strings.TrimSuffix(strings.Join(parts[2:], "/"), ".git"),
	}
	u.User = url.UserPassword("", o.GitHostAccessToken)
	client, err := defaultClientFactory(u.String())
	if err != nil {
		return fmt.Errorf("failed to create a client to access %q: %w", o.GitOpsRepoURL, err)
	}
	_, _, err = client.Repositories.Create(context.Background(), ri)
	return err
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
