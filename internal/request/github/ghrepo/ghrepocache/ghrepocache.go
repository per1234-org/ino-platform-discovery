// Package ghrepocache contains code related to the GitHub repositories data cache.
package ghrepocache

import (
	"github.com/per1234-org/ino-platform-discovery/internal/request/github/ghrepo/ghrepocache/ownercache"
	"github.com/per1234-org/ino-platform-discovery/internal/results/repo"
)

// Type is the type of the repositories data cache.
type Type map[string]ownercache.Type

// New returns a new cache object.
func New() Type {
	return make(Type)
}

// Get returns the repository data for the given owner and repository name.
func (cache Type) Get(owner string, name string) (repo.Type, bool) {
	ownerCache, cached := cache[owner]
	if !cached {
		return repo.Type{}, cached
	}

	repo, cached := ownerCache.Get(name)

	return repo, cached
}

// Set stores the repository data for the given owner and repository name.
func (cache Type) Set(owner string, name string, repoData repo.Type) {
	if cache[owner] == nil {
		cache[owner] = ownercache.New()
	}
	cache[owner].Set(name, repoData)
}
