// Package result contains code related to discovery result element.
package result

import (
	"github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result/content"
)

// Type is the type for the result data.
type Type struct {
	// Content is the type of content discovered in the result.
	Content content.Type
	// DefaultBranch is the name of the repository's default branch.
	DefaultBranch string
	// IndexPath is the path to the index file in the repository.
	IndexPath string
	// Owner is the username of the repository owner.
	Owner string
	// PlatformFilePath is the path to the platform data file in the repository discovered by the search.
	PlatformFilePath string
	// RepositoryData is the supplemental repository data.
	RepositoryData *repo.Type
	// RepositoryName is the name of the result's repository.
	RepositoryName string
	// RepositoryURL is the URL of the result's repository.
	RepositoryURL string
}

// Merge merges the data from two results.
func (result *Type) Merge(addition Type) {
	result.Content.Index = (result.Content.Index || addition.Content.Index)
	result.Content.Platform = (result.Content.Platform || addition.Content.Platform)

	if result.DefaultBranch == "" {
		result.DefaultBranch = addition.DefaultBranch
	}

	if result.IndexPath == "" {
		result.IndexPath = addition.IndexPath
	}

	if result.PlatformFilePath == "" {
		result.PlatformFilePath = addition.PlatformFilePath
	}

	if result.RepositoryURL == "" {
		result.RepositoryURL = addition.RepositoryURL
	}
}

// Same determines whether two results are for the same repository.
func (result *Type) Same(b Type) bool {
	return result.RepositoryURL == b.RepositoryURL
}
