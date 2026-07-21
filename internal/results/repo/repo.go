// Package repo contains code related to repository data.
package repo

// Type is the type of the repository data
type Type struct {
	// Ahead is whether a fork is ahead of its parent.
	Ahead bool
	// DefaultBranch is the name of the default branch.
	DefaultBranch string
	// Error is the error returned by the GitHub API request.
	Error error
	// Fork is whether the repository is a fork.
	Fork bool
	// ResolvedURL is the resolved repository URL.
	ResolvedURL string
}
