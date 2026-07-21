// Package request contains code related to requesting data that will be used for discovery.
package request

import (
	"github.com/google/go-github/v79/github"
	"github.com/per1234-org/ino-platform-discovery/internal/request/github/ghrepo"
	"github.com/per1234-org/ino-platform-discovery/internal/request/github/ghsearch"
	"github.com/per1234-org/ino-platform-discovery/internal/results"
	resultsrepo "github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result"
)

// Search searches for repositories that contain a package index and/or platform.
func Search(githubClient *github.Client) (results.Type, error) {
	return ghsearch.Search(githubClient)
}

// Supplement requests additional data and uses it to supplement the passed results.
func Supplement(client *github.Client, results *results.Type) error {
	repositoriesData := make(map[string](resultsrepo.Type))
	for resultIndex, result := range *results {
		// A repository may contain multiple platforms or indexes. In this case, it will be present multiple times in the
		// results. The repository data API request should only be made once for each repository in the results.
		repositoryData, populated := repositoriesData[result.RepositoryURL]
		if !populated {
			var err error
			repositoryData, err = repo(client, result)
			if err != nil {
				return err
			}

			repositoriesData[result.RepositoryURL] = repositoryData
		}

		result.RepositoryData = repositoryData

		(*results)[resultIndex] = result
	}

	return nil
}

// repo gets data for the given repository.
func repo(githubClient *github.Client, result result.Type) (resultsrepo.Type, error) {
	return ghrepo.Get(githubClient, result.Owner, result.RepositoryName)
}
