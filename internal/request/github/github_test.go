package github

import (
	"fmt"
	"testing"

	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit/github_primary_ratelimit"
	"github.com/stretchr/testify/assert"
)

// TestShouldRetry provides coverage for the `ShouldRetry` function.
func TestShouldRetry(t *testing.T) {
	var typedRateLimitReachedError *github_primary_ratelimit.RateLimitReachedError

	testTables := []struct {
		testName  string
		err       error
		assertion bool
	}{
		{
			"Return true when error is caused by rate limiting.",
			typedRateLimitReachedError,
			true,
		},
		{
			"Return true when error not caused by rate limiting.",
			fmt.Errorf("foo"),
			false,
		},
	}

	for _, testTable := range testTables {
		assert.Equal(t, testTable.assertion, ShouldRetry(testTable.err), testTable.testName)
	}
}
