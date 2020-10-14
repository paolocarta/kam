package pipelines

import (
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/driver/fake"
)

func TestBootstrapRepository(t *testing.T) {
	token := "this-is-a-test-token"
	fakeData := stubOutFactory(t, token)

	err := BootstrapRepository(&BootstrapOptions{
		GitOpsRepoURL:      "https://example.com/testing/test-repo.git",
		GitHostAccessToken: token,
	})
	assertNoError(t, err)

	want := []*scm.RepositoryInput{
		{
			Namespace:   "testing",
			Name:        "test-repo",
			Description: defaultRepoDescription,
			Private:     true,
		},
	}
	if diff := cmp.Diff(want, fakeData.CreateRepositories); diff != "" {
		t.Fatalf("BootstrapRepository failed:\n%s", diff)
	}
}

func TestRepoURL(t *testing.T) {
	urlTests := []struct {
		repoURL string
		wantURL string
	}{
		{"https://github.com/my-org/my-repo.git", "https://github.com"},
		{"https://gl.example.com/my-org/my-repo.git", "https://gl.example.com"},
	}

	for _, tt := range urlTests {
		t.Run(tt.repoURL, func(rt *testing.T) {
			u, err := repoURL(tt.repoURL)
			if err != nil {
				rt.Error(err)
				return
			}
			if u != tt.wantURL {
				rt.Errorf("got %q, want %q", u, tt.wantURL)
			}
		})
	}
}

func stubOutFactory(t *testing.T, authToken string) *fake.Data {
	f := defaultClientFactory
	t.Cleanup(func() {
		defaultClientFactory = f
	})

	client, data := fake.NewDefault()
	defaultClientFactory = func(repoURL string) (*scm.Client, error) {
		u, err := url.Parse(repoURL)
		if err != nil {
			return nil, err
		}
		want := ":" + authToken
		if a := u.User.String(); a != want {
			t.Fatalf("client failed auth: got %q, want %q", a, want)
		}
		return client, nil
	}
	return data
}
