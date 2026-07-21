// Package catalogentry contains code related to the individual catalog entries.
package catalogentry

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogcolumn"
	"github.com/per1234-org/ino-platform-discovery/internal/request/github/ghrepo"
	"github.com/sirupsen/logrus"
)

// New returns a new catalog entry object.
func New() []string {
	return make([]string, catalogcolumn.EnumEnd)
}

// IsDuplicate determines whether a candidate entry is a duplicate of an existing entry.
func IsDuplicate(incoming []string, existing []string) bool {
	/*
		The inoplatforms catalog entries contain data about the index associated with a platform in addition to the
		platform. Conversely, it is not possible for ino-platform-discovery to reliably determine such associations
		(presence in the same repository is not sufficient evidence because a repository may contain multiple platforms
		and/or indexes). So the "incoming" results only contain platform or index data alone. For this reason, checks for
		the "incoming" as a duplicate index and as a duplicate platform are performed separately.

		The branch name must be compared because, even though GitHub code search only covers the default branch, the catalog
		also contains items manually discovered in other branches. A catalog item may be present in the same repository as
		the result, but located in a non-default branch. In this case, the result is not a duplicate.
	*/

	// Check if it is a duplicate platform.
	if incoming[catalogcolumn.Repository] != "" {
		// Perform initial check with verbatim catalog entry to avoid unnecessary normalization requests.
		if isDuplicatePlatform(incoming, existing) {
			return true
		}

		resolvedExistingRepository, err := resolveURL(existing[catalogcolumn.Repository])
		if err == nil {
			if existing[catalogcolumn.Repository] != resolvedExistingRepository {
				existing[catalogcolumn.Repository] = resolvedExistingRepository
				if isDuplicatePlatform(incoming, existing) {
					return true
				}
			}
		}
	}

	// Check if it is a duplicate index.
	if incoming[catalogcolumn.PackageIndexRepository] != "" {
		if isDuplicateIndex(incoming, existing) {
			return true
		}

		resolvedExistingPackageIndexRepository, err := resolveURL(existing[catalogcolumn.PackageIndexRepository])
		if err == nil {
			if existing[catalogcolumn.PackageIndexRepository] != resolvedExistingPackageIndexRepository {
				existing[catalogcolumn.PackageIndexRepository] = resolvedExistingPackageIndexRepository
				if isDuplicateIndex(incoming, existing) {
					return true
				}
			}
		}
	}

	return false
}

// isDuplicateIndex determines whether a candidate index entry is a duplicate of an existing entry.
func isDuplicateIndex(incoming []string, existing []string) bool {
	incomingIndexFilename := urlFile(incoming[catalogcolumn.BoardsManagerURL])
	existingIndexFilename := urlFile(existing[catalogcolumn.BoardsManagerURL])

	return incoming[catalogcolumn.PackageIndexRepository] == existing[catalogcolumn.PackageIndexRepository] &&
		incoming[catalogcolumn.PackageIndexFolder] == existing[catalogcolumn.PackageIndexFolder] &&
		incoming[catalogcolumn.PackageIndexBranch] == existing[catalogcolumn.PackageIndexBranch] &&
		incomingIndexFilename == existingIndexFilename
}

// isDuplicatePlatform determines whether a candidate platform entry is a duplicate of an existing entry.
func isDuplicatePlatform(incoming []string, existing []string) bool {
	return incoming[catalogcolumn.Repository] == existing[catalogcolumn.Repository] &&
		incoming[catalogcolumn.RepositoryDataFolder] == existing[catalogcolumn.RepositoryDataFolder] &&
		incoming[catalogcolumn.BranchName] == existing[catalogcolumn.BranchName]
}

// resolveURL returns the given URL after any redirects are resolved.
func resolveURL(inputURL string) (string, error) {
	urlObject, err := url.Parse(inputURL)
	if err != nil {
		logrus.Errorf("Invalid URL: %s", inputURL)
		return inputURL, err
	}

	switch urlObject.Host {
	case "github.com":
		pathComponents := strings.Split(urlObject.Path, "/")
		if len(pathComponents) < 3 {
			logrus.Errorf("Incomplete GitHub repo URL: %s", inputURL)
			return inputURL, fmt.Errorf("incomplete GitHub repo URL: %s", inputURL)
		}

		owner := pathComponents[1]
		name := pathComponents[2]
		repoData, err := ghrepo.Get(owner, name)
		if err != nil {
			logrus.Errorf("GitHub API endpoint returned an error: %s", err.Error())
			return inputURL, err
		}

		return repoData.ResolvedURL, nil
	default:
		logrus.Tracef("URL not handled by resolver: %s", inputURL)
		return inputURL, nil
	}
}

// urlFile returns the last component of the path from the given URL.
func urlFile(rawURL string) string {
	urlObject, err := url.Parse(rawURL)
	if err != nil {
		logrus.Errorf("Invalid URL: %s", rawURL)
		return ""
	}

	if urlObject.Path == "" {
		// URL does not have a path component.
		logrus.Errorf("URL without path component: %s", rawURL)
		return ""
	}

	pathComponents := strings.Split(urlObject.Path, "/")

	return pathComponents[len(pathComponents)-1]
}
