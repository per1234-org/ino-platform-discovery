// Package github contains code for making GitHub API requests.
package github

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofri/go-github-pagination/githubpagination"
	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit"
	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit/github_primary_ratelimit"
	"github.com/gofri/go-github-ratelimit/v2/github_ratelimit/github_secondary_ratelimit"
	"github.com/google/go-github/v79/github"
	"github.com/sirupsen/logrus"
)

// NewClient creates a client for GitHub API requests.
func NewClient(githubToken string) (context.Context, *github.Client) {
	rateLimiterTransport := github_ratelimit.New(
		nil,
		github_primary_ratelimit.WithLimitDetectedCallback(
			func(ctx *github_primary_ratelimit.CallbackContext) {
				fmt.Printf(
					"\nReached GitHub API rate limit. Waiting until it resets at %s.\n",
					ctx.ResetTime.Format(time.DateTime),
				)
				logrus.Infof("Rate limit category: %s. HTTP status: %s", ctx.Category, ctx.Response.Status)

				// This would typically be done via the github_primary_ratelimit.WithSleepUntilReset() option. However, that
				// option overrides github_primary_ratelimit.WithLimitDetectedCallback so the sleep must be implemented in the
				// callback.
				time.Sleep(time.Until(*ctx.ResetTime))
				logrus.Info("Rate limit sleep finished.")
			},
		),
		// Unlike the primary rate limit handler, the secondary rate limit handler sleeps by default.
		github_secondary_ratelimit.WithLimitDetectedCallback(
			func(ctx *github_secondary_ratelimit.CallbackContext) {
				fmt.Printf(
					"\nReached secondary GitHub API rate limit. Waiting until it resets at %s.\n",
					ctx.ResetTime.Format(time.DateTime),
				)
			},
		),
	)

	/*
		Disable go-github's rate limit handling system so that it doesn't interfere with the superior rate limit handling by
		github_ratelimit.
		See:
		- https://github.com/gofri/go-github-ratelimit#usage-example-with-go-github
		- https://pkg.go.dev/github.com/google/go-github/v69/github#pkg-constants
	*/
	clientContext := context.WithValue(context.Background(), github.BypassRateLimitCheck, true)

	paginator := githubpagination.NewClient(rateLimiterTransport, githubpagination.WithPerPage(100))

	client := github.NewClient(paginator)
	// Token is needed to get a rate limiting allowance suitable for the full index, but some tests can run without.
	if githubToken != "" {
		logrus.Info("Using token authentication for GitHub API.")
		client = client.WithAuthToken(githubToken)
	}

	return clientContext, client
}

// ShouldRetry determines whether a failed request should be retried.
func ShouldRetry(err error) bool {
	// Source: https://github.com/gofri/go-github-ratelimit/blob/v2.0.0/github_ratelimit/primary_ratelimit_test.go#L211-L214
	var typedRateLimitReachedError *github_primary_ratelimit.RateLimitReachedError
	if errors.As(err, &typedRateLimitReachedError) {
		// If the rate limit handling system worked correctly, this would never occur. However, it doesn't, so
		// several requests may fail immediately after the sleep.
		logrus.Errorf("RateLimitReachedError during GitHub API request: %v", err)

		// Retry request.
		return true
	}

	return false
}
