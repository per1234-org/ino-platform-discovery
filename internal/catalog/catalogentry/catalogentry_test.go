package catalogentry

import (
	"testing"

	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogcolumn"
	"github.com/stretchr/testify/assert"
)

// TestNew provides coverage for the `New` function.
func TestNew(t *testing.T) {
	catalogEntry := New()
	assert.Equal(t, int(catalogcolumn.EnumEnd), len(catalogEntry))
}

// TestIsDuplicate provides coverage for the `IsDuplicate` function.
func TestIsDuplicate(t *testing.T) {
	existing := New()
	existing[catalogcolumn.Repository] = "existing Repository"
	existing[catalogcolumn.RepositoryDataFolder] = "existing RepositoryDataFolder"
	existing[catalogcolumn.BranchName] = "existing BranchName"

	incoming := New()
	incoming[catalogcolumn.Repository] = "existing Repository"
	incoming[catalogcolumn.RepositoryDataFolder] = "existing RepositoryDataFolder"
	incoming[catalogcolumn.BranchName] = "existing BranchName"

	assert.True(t, IsDuplicate(incoming, existing), "It should return true when incoming has duplicate platform data.")

	incoming[catalogcolumn.Repository] = "incoming Repository"

	assert.False(
		t,
		IsDuplicate(incoming, existing),
		"It should return false when incoming has unique platform data and empty index data.",
	)

	existing[catalogcolumn.PackageIndexRepository] = "existing PackageIndexRepository"
	existing[catalogcolumn.PackageIndexFolder] = "existing PackageIndexFolder"
	existing[catalogcolumn.PackageIndexBranch] = "existing PackageIndexBranch"

	incoming[catalogcolumn.PackageIndexRepository] = "existing PackageIndexRepository"
	incoming[catalogcolumn.PackageIndexFolder] = "existing PackageIndexFolder"
	incoming[catalogcolumn.PackageIndexBranch] = "existing PackageIndexBranch"

	assert.True(
		t,
		IsDuplicate(incoming, existing),
		"It should return true when incoming has unique platform data and duplicate index data.",
	)

	incoming[catalogcolumn.PackageIndexRepository] = "incoming PackageIndexRepository"

	assert.False(
		t,
		IsDuplicate(incoming, existing),
		"It should return false when incoming has unique platform data and unique index data.",
	)
}
