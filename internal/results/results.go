// Package results contains code related to the list of discoveries.
package results

import (
	"fmt"
	"slices"
	"strings"

	"github.com/per1234-org/ino-platform-discovery/internal/catalog"
	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogcolumn"
	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogentry"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result"
)

// Type is the type for the result data.
type Type []result.Type

// Deduplicate removes results that are already present in the catalog.
func (results *Type) Deduplicate(catalog catalog.Type) {
	deduplicated := slices.DeleteFunc(
		*results,
		func(result result.Type) bool {
			for _, catalogEntry := range catalog {
				if result.RepositoryURL == catalogEntry[catalogcolumn.Repository] || result.RepositoryURL == catalogEntry[catalogcolumn.PackageIndexRepository] {
					return true
				}
			}

			return false
		},
	)

	*results = deduplicated
}

// Merge merges the data from two sets of results.
func (results *Type) Merge(additions Type) {
	for _, addition := range additions {
		targetIndex := (*results).find(addition)
		if targetIndex == -1 {
			// Addition is not already present in target set.
			*results = append(*results, addition)
		} else {
			// An item matching addition is present in target set.
			(*results)[targetIndex].Merge(addition)
		}
	}
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
		indexBoardsManagerURL := ""
		indexBranch := ""
		indexFolder := ""
		indexRepository := ""
		if result.Content.Index {
			// E.g., https://raw.githubusercontent.com/damellis/attiny/refs/heads/ide-1.6.x-boards-manager/package_damellis_attiny_index.json
			indexBoardsManagerURL = fmt.Sprintf(
				"https://raw.githubusercontent.com/%s/%s/refs/heads/%s/%s",
				result.Owner,
				result.RepositoryName,
				result.DefaultBranch,
				result.IndexPath,
			)

			indexBranch = result.DefaultBranch
			indexFolder = toCatalogPath(result.IndexPath)
			indexRepository = result.RepositoryURL
		}

		platformBranch := ""
		platformFolder := ""
		platformRepository := ""
		if result.Content.Platform {
			platformBranch = result.DefaultBranch
			platformFolder = toCatalogPath(result.PlatformFilePath)
			platformRepository = result.RepositoryURL
		}

		catalogEntry := catalogentry.New()

		catalogEntry[catalogcolumn.BoardsManagerURL] = indexBoardsManagerURL
		catalogEntry[catalogcolumn.PackageIndexBranch] = indexBranch
		catalogEntry[catalogcolumn.PackageIndexFolder] = indexFolder
		catalogEntry[catalogcolumn.PackageIndexRepository] = indexRepository
		catalogEntry[catalogcolumn.BranchName] = platformBranch
		catalogEntry[catalogcolumn.RepositoryDataFolder] = platformFolder
		catalogEntry[catalogcolumn.Repository] = platformRepository

		catalog = append(catalog, catalogEntry)
	}

	return catalog
}

// find finds the index of the item matching the given result in the given set of results.
// If the item is not found, -1 is returned.
func (results Type) find(query result.Type) int {
	for index, result := range results {
		if result.Same(query) {
			return index
		}
	}

	return -1
}

// toCatalogPath converts the value of a result path field to the appropriate format for use in the associated catalog
// field.
func toCatalogPath(resultPath string) string {
	resultPathComponents := strings.Split(resultPath, "/")
	if len(resultPathComponents) == 1 {
		// Index is in the root of the repository.
		return "/"
	}

	resultPathParentComponents := resultPathComponents[:len(resultPathComponents)-1]
	resultPathParentString := strings.Join(resultPathParentComponents, "/")
	// `resultPath` is the relative path of the file in the repository (e.g., `foo/bar/package_foo_index.json`).
	// By convention, the catalog entry has a leading and trailing separator (e.g., `/foo/bar/`).
	return fmt.Sprintf("/%s/", resultPathParentString)
}
