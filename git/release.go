package git

import (
	"context"

	"github.com/google/go-github/v32/github"
)

// GetRepoStablRelease list stable release for given repo
func GetRepoStablRelease(user string, repo string) ([]*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(context.Background(), user, repo, nil)
	if err != nil {
		return nil, err
	}

	var stableReleases []*github.RepositoryRelease
	for _, release := range releases {
		if !release.GetPrerelease() {
			stableReleases = append(stableReleases, release)
		}
	}

	return stableReleases, nil
}
