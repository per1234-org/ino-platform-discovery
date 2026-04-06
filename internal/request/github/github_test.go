package github

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-github/v79/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHandleRateLimiting provides coverage for the `HandleRateLimiting` function.
func TestHandleRateLimiting(t *testing.T) {
	primaryRateLimitErr := github.RateLimitError{
		Rate: github.Rate{
			Reset: github.Timestamp{
				Time: time.Now(),
			},
		},
	}

	assert.Nil(t, HandleRateLimiting(&primaryRateLimitErr), "Return nil when error is caused by primary rate limiting.")

	retryAfter, err := time.ParseDuration("0s")
	require.NoError(t, err)
	secondaryRateLimitErr := github.AbuseRateLimitError{
		RetryAfter: &retryAfter,
	}

	assert.Nil(t, HandleRateLimiting(&secondaryRateLimitErr), "Return nil when error is caused by secondary rate limiting.")

	otherError := errors.New("foo")

	assert.Error(t, HandleRateLimiting(otherError), "Return error when error is caused by something other than rate limiting.")
}
