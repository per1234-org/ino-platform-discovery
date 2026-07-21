// Package ghrepocache contains code related to the GitHub repositories data cache.
package ghrepocache

import "github.com/per1234-org/ino-platform-discovery/internal/request/github/ghrepo/ghrepocache/ownercache"

// Type is the type of the repositories data cache.
type Type map[string]ownercache.Type

// New returns a new cache object.
func New() Type {
	return make(Type)
}
