// Package ownercache contains code related to the per-owner repositories data cache.
package ownercache

import "github.com/per1234-org/ino-platform-discovery/internal/results/repo"

// Type is the type of the per-owner cache data.
type Type map[string]repo.Type

// New returns a new cache object.
func New() Type {
	return make(Type)
}
