package request

import (
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
	githubContext, githubClient := github.NewClient("")

	targetResults := results.Type{
		result.Type{
			Owner:          "per1234-org",
			RepositoryName: "Enterprise",
		},
	}
	err := Supplement(githubContext, githubClient, &targetResults)
	require.NoError(t, err)

	repo := resultsrepo.Type{
		Ahead:         false,
		DefaultBranch: "master",
		Fork:          true,
	}

	assertion := results.Type{
		result.Type{
			Owner:          "per1234-org",
			RepositoryData: &repo,
			RepositoryName: "Enterprise",
		},
	}

	assert.Equal(t, assertion, targetResults, "Repo data as expected.")
}

// TestLoad provides coverage for the `repo` function.
func Test_repo(t *testing.T) {
	githubContext, githubClient := github.NewClient("")

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

	assert.Equal(t, assertion, *repo, "Repo data as expected.")
}
