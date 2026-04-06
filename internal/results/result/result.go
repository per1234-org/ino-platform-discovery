// Package result contains code related to discovery result element.
package result

import (
	"fmt"
	"strings"

	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogcolumn"
	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogentry"
	"github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result/content"
)

// Type is the type for the result data.
type Type struct {
	// Content is the type of content discovered in the result.
	Content content.Type
	// Owner is the username of the repository owner.
	Owner string
	// Path is the path of the discovered file in the repository.
	Path string
	// RepositoryData is the supplemental repository data.
	RepositoryData repo.Type
	// RepositoryName is the name of the result's repository.
	RepositoryName string
	// RepositoryURL is the URL of the result's repository.
	RepositoryURL string
}

// ToCatalogEntry returns the given result in the catalog entry data format.
func (result Type) ToCatalogEntry() []string {
	indexBoardsManagerURL := ""
	indexBranch := ""
	indexFolder := ""
	indexRepository := ""

	platformBranch := ""
	platformFolder := ""
	platformRepository := ""

	switch result.Content {
	case content.Index:
		// E.g., https://raw.githubusercontent.com/damellis/attiny/refs/heads/ide-1.6.x-boards-manager/package_damellis_attiny_index.json
		indexBoardsManagerURL = fmt.Sprintf(
			"https://raw.githubusercontent.com/%s/%s/refs/heads/%s/%s",
			result.Owner,
			result.RepositoryName,
			result.RepositoryData.DefaultBranch,
			result.Path,
		)

		indexBranch = result.RepositoryData.DefaultBranch
		indexFolder = toCatalogPath(result.Path)
		indexRepository = result.RepositoryURL

	case content.Platform:
		platformBranch = result.RepositoryData.DefaultBranch
		platformFolder = toCatalogPath(result.Path)
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

	return catalogEntry
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
