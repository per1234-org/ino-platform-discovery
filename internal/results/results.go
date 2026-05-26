// Package results contains code related to the list of discoveries.
package results

import (
	"slices"

	"github.com/per1234-org/ino-platform-discovery/internal/catalog"
	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogentry"
	"github.com/per1234-org/ino-platform-discovery/internal/exclusions"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result"
)

// Type is the type for the result data.
type Type []result.Type

// Deduplicate removes results that are already present in the catalog.
func (results *Type) Deduplicate(catalog catalog.Type) {
	deduplicated := slices.DeleteFunc(
		*results,
		func(result result.Type) bool {
			resultEntry := result.ToCatalogEntry()
			for _, catalogEntry := range catalog {
				if catalogentry.IsDuplicate(resultEntry, catalogEntry) {
					// Result is duplicate, delete.
					return true
				}
			}

			// Result is novel, retain.
			return false
		},
	)

	*results = deduplicated
}

// Exclude removes excluded results.
func (results *Type) Exclude(exclusions exclusions.Type) {
	included := slices.DeleteFunc(
		*results,
		func(result result.Type) bool {
			for _, exclusion := range exclusions {
				if exclusion.Match(result) {
					// Result is to be excluded, delete.
					return true
				}
			}

			// Result is not to be excluded, retain.
			return false
		},
	)

	*results = included
}

// Filter removes results determined to not be valid discoveries.
func (results *Type) Filter() {
	filtered := slices.DeleteFunc(
		*results,
		func(result result.Type) bool {
			if result.RepositoryData.Fork && !result.RepositoryData.Ahead {
				// Filter out forks that are not ahead of the parent repo.
				return true
			}

			// Retain the result.
			return false
		},
	)

	*results = filtered
}

// ToCatalog returns the given results in the catalog data format.
func (results Type) ToCatalog() catalog.Type {
	catalog := catalog.Type{}
	for _, result := range results {
		catalog = append(catalog, result.ToCatalogEntry())
	}

	return catalog
}
