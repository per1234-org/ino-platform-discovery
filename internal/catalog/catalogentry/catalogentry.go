// Package catalogentry contains code related to the individual catalog entries.
package catalogentry

import (
	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogcolumn"
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
	if incoming[catalogcolumn.Repository] != "" &&
		incoming[catalogcolumn.Repository] == existing[catalogcolumn.Repository] &&
		incoming[catalogcolumn.RepositoryDataFolder] == existing[catalogcolumn.RepositoryDataFolder] &&
		incoming[catalogcolumn.BranchName] == existing[catalogcolumn.BranchName] {
		return true
	}

	// Check if it is a duplicate index.
	if incoming[catalogcolumn.PackageIndexRepository] != "" &&
		incoming[catalogcolumn.PackageIndexRepository] == existing[catalogcolumn.PackageIndexRepository] &&
		incoming[catalogcolumn.PackageIndexFolder] == existing[catalogcolumn.PackageIndexFolder] &&
		incoming[catalogcolumn.PackageIndexBranch] == existing[catalogcolumn.PackageIndexBranch] {
		return true
	}

	return false
}
