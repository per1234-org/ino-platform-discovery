// Package repo contains code related to repository data.
package repo

// Type is the type of the repository data
type Type struct {
	// Ahead is whether a fork is ahead of its parent.
	Ahead bool
	// DefaultBranch is the name of the default branch.
	DefaultBranch string
	// Fork is whether the repository is a fork.
	Fork bool
}
