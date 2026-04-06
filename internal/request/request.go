// Package request contains code related to requesting data that will be used for discovery.
package request

import (
	"context"

	"github.com/google/go-github/v79/github"
	"github.com/per1234-org/ino-platform-discovery/internal/request/github/ghrepo"
	"github.com/per1234-org/ino-platform-discovery/internal/request/github/ghsearch"
	"github.com/per1234-org/ino-platform-discovery/internal/results"
	resultsrepo "github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result"
)

// Search searches for repositories that contain a package index and/or platform.
func Search(githubContext context.Context, githubClient *github.Client) (results.Type, error) {
	return ghsearch.Search(githubContext, githubClient)
}

// Supplement requests additional data and uses it to supplement the passed results.
func Supplement(clientContext context.Context, client *github.Client, results *results.Type) error {
	for resultIndex, result := range *results {
		var err error
		result.RepositoryData, err = repo(clientContext, client, result)
		if err != nil {
			return err
		}

		(*results)[resultIndex] = result
	}

	return nil
}

// repo gets data for the given repository.
func repo(githubContext context.Context, githubClient *github.Client, result result.Type) (*resultsrepo.Type, error) {
	return ghrepo.Get(githubContext, githubClient, result.Owner, result.RepositoryName)
}
