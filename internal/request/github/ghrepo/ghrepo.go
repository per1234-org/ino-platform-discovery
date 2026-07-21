// Package ghrepo contains code for making requests to the `/repos/{owner}/{repo}` endpoint of the GitHub API.
package ghrepo

import (
	"context"
	"fmt"

	gogithub "github.com/google/go-github/v79/github"
	"github.com/per1234-org/ino-platform-discovery/internal/request/clients"
	"github.com/per1234-org/ino-platform-discovery/internal/request/github"
	"github.com/per1234-org/ino-platform-discovery/internal/request/github/ghrepo/ghrepocache"
	"github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/sirupsen/logrus"
)

var repoCache ghrepocache.Type

// Get makes a request to the `/repos/{owner}/{repo}` endpoint of the GitHub API and returns data extracted from the response.
func Get(owner string, name string) (repo.Type, error) {
	if repoCache == nil {
		repoCache = ghrepocache.New()
	}

	repo := repo.Type{}

	if repo, cached := repoCache.Get(owner, name); cached {
		// Use cached data instead of performing redundant request.
		return repo, repo.Error
	}

	var githubResponse *gogithub.Repository
	for {
		logrus.Tracef("Making GitHub API /repos/{owner}/{repo} endpoint request for %s/%s", owner, name)
		var err error
		githubResponse, _, err = clients.Clients.GitHub.Repositories.Get(context.Background(), owner, name)
		repo.Error = err

		if err != nil {
			err := github.HandleRateLimiting(err)
			if err != nil {
				// Error is not recoverable.
				repoCache.Set(owner, name, repo)
				return repo, err
			}

			// Retry request.
			continue
		}

		// Request was successful.
		break
	}

	repo.DefaultBranch = *githubResponse.DefaultBranch
	repo.Fork = *githubResponse.Fork
	repo.ResolvedURL = *githubResponse.HTMLURL

	if repo.Fork {
		ahead, err := ahead(githubResponse)
		if err != nil {
			return repo, err
		}

		repo.Ahead = ahead
	}

	repoCache.Set(owner, name, repo)
	return repo, nil
}

// ahead makes a request to the `/repos/{owner}/{repo}/compare/{basehead}` endpoint of the GitHub API for the subject
// repo and its parent, then returns whether the subject repo is "ahead" of the parent.
func ahead(getRepoResponse *gogithub.Repository) (bool, error) {
	ahead := false
	for {
		logrus.Tracef("Making GitHub API /repos/{owner}/{repo}/compare/{basehead} endpoint request for %s", *getRepoResponse.FullName)

		// E.g., https://api.github.com/repos/arduino/arduino-esp32/compare/espressif:master...arduino:master
		base := fmt.Sprintf("%s:%s", *getRepoResponse.Parent.Owner.Login, *getRepoResponse.Parent.DefaultBranch)
		head := fmt.Sprintf("%s:%s", *getRepoResponse.Owner.Login, *getRepoResponse.DefaultBranch)
		githubResponse, _, err := clients.Clients.GitHub.Repositories.CompareCommits(
			context.Background(),
			*getRepoResponse.Owner.Login,
			*getRepoResponse.Name,
			base,
			head,
			nil,
		)

		if err != nil {
			err := github.HandleRateLimiting(err)
			if err != nil {
				// Error is not recoverable.
				return ahead, err
			}

			// Retry request.
			continue
		}

		ahead = *githubResponse.AheadBy > 0

		// Request was successful.
		break
	}

	return ahead, nil
}
