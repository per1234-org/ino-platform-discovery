package request

import (
	"os"
	"testing"

	"github.com/per1234-org/ino-platform-discovery/internal/request/github"
	"github.com/per1234-org/ino-platform-discovery/internal/results"
	resultsrepo "github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSupplement provides coverage for the `Supplement` function.
func TestSupplement(t *testing.T) {
	if os.Getenv("GITHUB_TOKEN") == "" {
		// It is possible to run the test without a token. However, this makes it flakey due to heavy rate limiting.
		t.Skip("Required GITHUB_TOKEN environment variable not defined.")
	}

	githubContext, githubClient := github.NewClient(os.Getenv("GITHUB_TOKEN"))

	targetResults := results.Type{
		result.Type{
			Owner:          "per1234-org",
			RepositoryName: "Enterprise",
			RepositoryURL:  "github.com/per1234-org/Enterprise",
		},
	}
	err := Supplement(githubContext, githubClient, &targetResults)
	require.NoError(t, err)

	assertion := results.Type{
		result.Type{
			Owner: "per1234-org",
			RepositoryData: resultsrepo.Type{
				Ahead:         false,
				DefaultBranch: "master",
				Fork:          true,
			},
			RepositoryName: "Enterprise",
			RepositoryURL:  "github.com/per1234-org/Enterprise",
		},
	}

	assert.Equal(t, assertion, targetResults, "Repo data as expected.")
}

// Test_repo provides coverage for the `repo` function.
func Test_repo(t *testing.T) {
	if os.Getenv("GITHUB_TOKEN") == "" {
		// It is possible to run the test without a token. However, this makes it flakey due to heavy rate limiting.
		t.Skip("Required GITHUB_TOKEN environment variable not defined.")
	}

	githubContext, githubClient := github.NewClient(os.Getenv("GITHUB_TOKEN"))

	result := result.Type{
		Owner:          "per1234-org",
		RepositoryName: "Enterprise",
	}
	repo, err := repo(githubContext, githubClient, result)
	require.NoError(t, err)

	assertion := resultsrepo.Type{
		Ahead:         false,
		DefaultBranch: "master",
		Fork:          true,
	}

	assert.Equal(t, assertion, repo, "Repo data as expected.")
}
