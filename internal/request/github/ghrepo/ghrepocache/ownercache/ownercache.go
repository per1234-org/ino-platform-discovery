// Package ownercache contains code related to the per-owner repositories data cache.
package ownercache

import "github.com/per1234-org/ino-platform-discovery/internal/results/repo"

// Type is the type of the per-owner cache data.
type Type map[string]repo.Type

// New returns a new cache object.
func New() Type {
	return make(Type)
}

// Get returns the repository data for the given repository name.
func (cache Type) Get(name string) (repo.Type, bool) {
	repo, cached := cache[name]

	return repo, cached
}

// Set stores the repository data for the given repository name.
func (cache Type) Set(name string, repoData repo.Type) {
	cache[name] = repoData
}
