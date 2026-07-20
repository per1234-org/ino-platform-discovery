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
	testTables := []struct {
		testMessage                    string
		incomingRepository             string
		incomingBoardsManagerURL       string
		incomingRepositoryDataFolder   string
		incomingBranchName             string
		incomingPackageIndexRepository string
		incomingPackageIndexFolder     string
		incomingPackageIndexBranch     string
		existingRepository             string
		existingBoardsManagerURL       string
		existingRepositoryDataFolder   string
		existingBranchName             string
		existingPackageIndexRepository string
		existingPackageIndexFolder     string
		existingPackageIndexBranch     string
		expected                       bool
	}{
		{
			testMessage:                  "It should return true when incoming is duplicate platform.",
			incomingRepository:           "https://github.com/foo-owner/foo-repo",
			incomingRepositoryDataFolder: "foo-folder",
			incomingBranchName:           "foo-branch",
			existingRepository:           "https://github.com/foo-owner/foo-repo",
			existingRepositoryDataFolder: "foo-folder",
			existingBranchName:           "foo-branch",
			expected:                     true,
		},
		{
			testMessage:                  "It should return false when incoming is novel platform.",
			incomingRepository:           "https://github.com/bar-owner/bar-repo",
			incomingRepositoryDataFolder: "bar-folder",
			incomingBranchName:           "bar-branch",
			existingRepository:           "https://github.com/foo-owner/foo-repo",
			existingRepositoryDataFolder: "foo-folder",
			existingBranchName:           "foo-branch",
			expected:                     false,
		},
		{
			testMessage:                    "It should return true when incoming is duplicate index.",
			incomingPackageIndexRepository: "https://github.com/foo-owner/foo-repo",
			incomingBoardsManagerURL:       "https://raw.githubusercontent.com/foo-owner/foo-repo/refs/heads/main/package_foo_index.json",
			incomingPackageIndexFolder:     "foo-folder",
			incomingPackageIndexBranch:     "foo-branch",
			existingPackageIndexRepository: "https://github.com/foo-owner/foo-repo",
			existingBoardsManagerURL:       "https://raw.githubusercontent.com/foo-owner/foo-repo/refs/heads/main/package_foo_index.json",
			existingPackageIndexFolder:     "foo-folder",
			existingPackageIndexBranch:     "foo-branch",
			expected:                       true,
		},
		{
			testMessage:                    "It should return false when incoming is novel index.",
			incomingPackageIndexRepository: "https://github.com/bar-owner/bar-repo",
			incomingBoardsManagerURL:       "https://raw.githubusercontent.com/bar-owner/bar-repo/refs/heads/main/package_bar_index.json",
			incomingPackageIndexFolder:     "bar-folder",
			incomingPackageIndexBranch:     "bar-branch",
			existingPackageIndexRepository: "https://github.com/foo-owner/foo-repo",
			existingBoardsManagerURL:       "https://raw.githubusercontent.com/foo-owner/foo-repo/refs/heads/main/package_foo_index.json",
			existingPackageIndexFolder:     "foo-folder",
			existingPackageIndexBranch:     "foo-branch",
			expected:                       false,
		},
	}

	for _, testTable := range testTables {
		incoming := New()
		incoming[catalogcolumn.Repository] = testTable.incomingRepository
		incoming[catalogcolumn.RepositoryDataFolder] = testTable.incomingRepositoryDataFolder
		incoming[catalogcolumn.BranchName] = testTable.incomingBranchName
		incoming[catalogcolumn.PackageIndexRepository] = testTable.incomingPackageIndexRepository
		incoming[catalogcolumn.PackageIndexFolder] = testTable.incomingPackageIndexFolder
		incoming[catalogcolumn.PackageIndexBranch] = testTable.incomingPackageIndexBranch

		existing := New()
		existing[catalogcolumn.Repository] = testTable.existingRepository
		existing[catalogcolumn.RepositoryDataFolder] = testTable.existingRepositoryDataFolder
		existing[catalogcolumn.BranchName] = testTable.existingBranchName
		existing[catalogcolumn.PackageIndexRepository] = testTable.existingPackageIndexRepository
		existing[catalogcolumn.PackageIndexFolder] = testTable.existingPackageIndexFolder
		existing[catalogcolumn.PackageIndexBranch] = testTable.existingPackageIndexBranch

		assert.Equal(t, testTable.expected, IsDuplicate(incoming, existing), testTable.testMessage)
	}
}

// Test_urlFile provides coverage for the `urlFile` function.
func Test_urlFile(t *testing.T) {
	assert.Equal(t, urlFile(":"), "", "It should return empty string when URL is invalid.")

	assert.Equal(
		t,
		urlFile("https://example.com"),
		"",
		"It should return empty string when URL does not have path component.",
	)

	assert.Equal(
		t,
		urlFile("https://example.com/"),
		"",
		"It should return empty string when URL does not have file component.",
	)

	assert.Equal(
		t,
		urlFile("https://example.com/foo.bar"),
		"foo.bar",
		"It should return the filename when URL has a file component in root.",
	)

	assert.Equal(
		t,
		urlFile("https://example.com/foo/bar/baz.qux"),
		"baz.qux",
		"It should return the filename when URL has a file component under subfolders.",
	)
}
