// Package github contains code for making GitHub API requests.
package github

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/go-github/v79/github"
	"github.com/sirupsen/logrus"
)

// HandleRateLimiting handles API request failures caused by rate limiting.
func HandleRateLimiting(err error) error {
	var primaryRateLimitErr *github.RateLimitError
	if errors.As(err, &primaryRateLimitErr) {
		// If you retry at the time specified by GitHub, the requests continue to fail due to primary rate limit for some
		// additional seconds.
		resetOffset := time.Second * 3

		resetTime := primaryRateLimitErr.Rate.Reset.Add(resetOffset)
		fmt.Printf("Reached primary GitHub API rate limit. Waiting until it resets at %s.\n", resetTime)
		time.Sleep(time.Until(resetTime))

		return nil
	}

	var secondaryRateLimitErr *github.AbuseRateLimitError
	if errors.As(err, &secondaryRateLimitErr) {
		fmt.Printf(
			"Reached secondary GitHub API rate limit. Waiting until it resets at %s.\n",
			time.Now().Add(*secondaryRateLimitErr.RetryAfter),
		)
		time.Sleep(*secondaryRateLimitErr.RetryAfter)

		return nil
	}

	// Error was caused by something other than rate limiting.
	return err
}

// NewClient creates a client for GitHub API requests.
func NewClient(githubToken string) (context.Context, *github.Client) {
	client := github.NewClient(nil)

	// Token is needed to get a rate limiting allowance suitable for the full index, but some tests can run without.
	if githubToken != "" {
		logrus.Info("Using token authentication for GitHub API.")
		client = client.WithAuthToken(githubToken)
	}

	return context.Background(), client
}
