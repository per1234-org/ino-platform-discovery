package ghrepo

import (
	"os"
	"testing"

	gogithub "github.com/google/go-github/v79/github"
	"github.com/per1234-org/ino-platform-discovery/internal/request/clients"
	"github.com/per1234-org/ino-platform-discovery/internal/request/github"
	resultsrepo "github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGet provides coverage for the `Get` function.
func TestGet(t *testing.T) {
	if os.Getenv("GITHUB_TOKEN") == "" {
		// It is possible to run the test without a token. However, this makes it flakey due to heavy rate limiting.
		t.Skip("Required GITHUB_TOKEN environment variable not defined.")
	}

	clients.Clients.GitHub = github.NewClient(os.Getenv("GITHUB_TOKEN"))

	repo, err := Get("per1234-org", "Enterprise")
	require.NoError(t, err)

	assertion := resultsrepo.Type{
		Ahead:         false,
		DefaultBranch: "master",
		Error:         nil,
		Fork:          true,
		ResolvedURL:   "https://github.com/per1234-org/Enterprise",
	}

	assert.Equal(t, assertion, repo, "Repo data as expected.")
}

// Test_ahead provides coverage for the `ahead` function.
func Test_ahead(t *testing.T) {
	if os.Getenv("GITHUB_TOKEN") == "" {
		// It is possible to run the test without a token. However, this makes it flakey due to heavy rate limiting.
		t.Skip("Required GITHUB_TOKEN environment variable not defined.")
	}

	clients.Clients.GitHub = github.NewClient(os.Getenv("GITHUB_TOKEN"))

	getRepoResponse := gogithub.Repository{
		DefaultBranch: new("master"),
		FullName:      new("per1234-org/Enterprise"),
		Name:          new("Enterprise"),
		Owner: &gogithub.User{
			Login: new("per1234-org"),
		},
		Parent: &gogithub.Repository{
			Owner: &gogithub.User{
				Login: new("delta-G"),
			},
			DefaultBranch: new("master"),
		},
	}

	isAhead, err := ahead(&getRepoResponse)
	require.NoError(t, err)

	assert.False(t, isAhead, "Fork without \"ahead\" status is not considered ahead.")

	getRepoResponse = gogithub.Repository{
		DefaultBranch: new("master"),
		FullName:      new("bhunting/libpololu-avr"),
		Name:          new("libpololu-avr"),
		Owner: &gogithub.User{
			Login: new("bhunting"),
		},
		Parent: &gogithub.Repository{
			Owner: &gogithub.User{
				Login: new("pololu"),
			},
			DefaultBranch: new("master"),
		},
	}

	isAhead, err = ahead(&getRepoResponse)
	require.NoError(t, err)

	assert.False(t, isAhead, "Fork with \"ahead\" status, but null diff is not considered ahead.")

	getRepoResponse = gogithub.Repository{
		DefaultBranch: new("main"),
		FullName:      new("arduino/ArduinoCore-zephyr"),
		Name:          new("ArduinoCore-zephyr"),
		Owner: &gogithub.User{
			Login: new("arduino"),
		},
		Parent: &gogithub.Repository{
			Owner: &gogithub.User{
				Login: new("zephyrproject-rtos"),
			},
			DefaultBranch: new("main"),
		},
	}

	isAhead, err = ahead(&getRepoResponse)
	require.NoError(t, err)

	assert.True(t, isAhead, "Fork with \"ahead\" status, and diff is considered ahead.")
}
