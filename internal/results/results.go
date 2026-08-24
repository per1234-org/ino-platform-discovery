// Package results contains code related to the list of discoveries.
package results

import (
	"slices"

	"github.com/per1234-org/ino-platform-discovery/internal/catalog"
	"github.com/per1234-org/ino-platform-discovery/internal/data"
	"github.com/per1234-org/ino-platform-discovery/internal/exclusions"
	"github.com/per1234-org/ino-platform-discovery/internal/feedback"
	"github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result/content"
)

// Type is the type for the result data.
type Type []result.Type

// Deduplicate removes results that are already present in the catalog.
func (results *Type) Deduplicate(catalog catalog.Type) {
	deduplicated := Type{}

	resultCount := len(*results)
	for resultIndex, candidateResult := range *results {
		feedback.Progress(resultIndex+1, resultCount)

		if !candidateResult.IsDuplicate(catalog) {
			deduplicated = append(deduplicated, candidateResult)
		}
	}

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

// FilterSupplemented removes results determined to not be valid discoveries based on the supplementary data.
func (results *Type) FilterSupplemented() {
	filtered := slices.DeleteFunc(
		*results,
		func(result result.Type) bool {
			if result.RepositoryData == (repo.Type{}) {
				panic("result has not been supplemented")
			}

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

// Prefilter removes results determined to not be valid discoveries based on the data returned by the search.
func (results *Type) Prefilter() {
	filtered := slices.DeleteFunc(
		*results,
		func(result result.Type) bool {
			/*
				The platform search query uses the `filename` qualifier to search for files named `boards.txt`. This qualifier
				also matches against filenames that are a substring match (e.g., a file named `foo.boards.txt` matches against
				the `filename:boards.txt` query). Any result that doesn't contain a file named exactly `boards.txt` is not valid
				and thus must be filtered out.

				An equivalent check is done on the filename of index results, but that is done in `ghsearch.indexes` instead of
				here in order to avoid unnecessary API requests.
			*/
			if result.Content == content.Platform && result.Filename != data.PlatformIndicatorFile {
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
