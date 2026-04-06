package ghrepo

import (
	"testing"

	"github.com/per1234-org/ino-platform-discovery/internal/request/github"
	resultsrepo "github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGet provides coverage for the `Get` function.
func TestGet(t *testing.T) {
	githubContext, githubClient := github.NewClient("")

	repo, err := Get(githubContext, githubClient, "per1234-org", "Enterprise")
	require.NoError(t, err)

	assertion := resultsrepo.Type{
		Ahead:         false,
		DefaultBranch: "master",
		Fork:          true,
	}

	assert.Equal(t, assertion, *repo, "Repo data as expected.")
}
