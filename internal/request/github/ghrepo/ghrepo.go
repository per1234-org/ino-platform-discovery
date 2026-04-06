// Package ghrepo contains code for making requests to the `/repos/{owner}/{repo}` endpoint of the GitHub API.
package ghrepo

import (
	"context"
	"fmt"

	gogithub "github.com/google/go-github/v79/github"
	"github.com/per1234-org/ino-platform-discovery/internal/request/github"
	"github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/sirupsen/logrus"
)

// Get makes a request to the `/repos/{owner}/{repo}` endpoint of the GitHub API and returns data extracted from the response.
func Get(clientContext context.Context, client *gogithub.Client, owner string, name string) (repo.Type, error) {
	repo := repo.Type{}

	var githubResponse *gogithub.Repository
	for {
		logrus.Tracef("Making GitHub API /repos/{owner}/{repo} endpoint request for %s/%s", owner, name)
		var err error
		githubResponse, _, err = client.Repositories.Get(clientContext, owner, name)

		if err != nil {
			err := github.HandleRateLimiting(err)
			if err != nil {
				// Error is not recoverable.
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

	if repo.Fork {
		ahead, err := ahead(clientContext, client, githubResponse)
		if err != nil {
			return repo, err
		}

		repo.Ahead = ahead
	}

	return repo, nil
}

// ahead makes a request to the `/repos/{owner}/{repo}/compare/{basehead}` endpoint of the GitHub API for the subject
// repo and its parent, then returns whether the subject repo is "ahead" of the parent.
func ahead(clientContext context.Context, client *gogithub.Client, getRepoResponse *gogithub.Repository) (bool, error) {
	ahead := false
	for {
		logrus.Tracef("Making GitHub API /repos/{owner}/{repo}/compare/{basehead} endpoint request for %s", *getRepoResponse.FullName)

		// E.g., https://api.github.com/repos/arduino/arduino-esp32/compare/espressif:master...arduino:master
		base := fmt.Sprintf("%s:%s", *getRepoResponse.Parent.Owner.Login, *getRepoResponse.Parent.DefaultBranch)
		head := fmt.Sprintf("%s:%s", *getRepoResponse.Owner.Login, *getRepoResponse.DefaultBranch)
		githubResponse, _, err := client.Repositories.CompareCommits(
			clientContext,
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
