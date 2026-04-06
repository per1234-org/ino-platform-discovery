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
func Get(clientContext context.Context, client *gogithub.Client, owner string, name string) (*repo.Type, error) {
	repo := repo.Type{}

	doRequest := true
	for doRequest {
		logrus.Tracef("Making GitHub API /repos/{owner}/{repo} endpoint request for %s/%s", owner, name)
		doRequest = false
		githubResponse, _, err := client.Repositories.Get(clientContext, owner, name)

		if err != nil {
			if github.ShouldRetry(err) {
				// Retry request.
				doRequest = true
				continue
			}

			// Error is not recoverable.
			return nil, err
		}

		repo.DefaultBranch = *githubResponse.DefaultBranch
		repo.Fork = *githubResponse.Fork

		if repo.Fork {
			repo.Ahead, err = ahead(clientContext, client, githubResponse)
			if err != nil {
				return nil, err
			}
		}
	}

	return &repo, nil
}

// ahead makes a request to the `/repos/{owner}/{repo}/compare/{basehead}` endpoint of the GitHub API for the subject
// repo and its parent, then returns whether the subject repo is "ahead" of the parent.
func ahead(clientContext context.Context, client *gogithub.Client, getRepoResponse *gogithub.Repository) (bool, error) {
	ahead := false
	doRequest := true
	for doRequest {
		logrus.Tracef("Making GitHub API /repos/{owner}/{repo}/compare/{basehead} endpoint request for %s", *getRepoResponse.FullName)
		doRequest = false

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
			if github.ShouldRetry(err) {
				// Retry request.
				doRequest = true
				continue
			}

			// Error is not recoverable.
			return ahead, err
		}

		ahead = *githubResponse.AheadBy > 0
	}

	return ahead, nil
}
