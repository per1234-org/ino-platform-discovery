package catalogentry

import (
	"os"
	"testing"

	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogcolumn"
	"github.com/per1234-org/ino-platform-discovery/internal/request/clients"
	"github.com/per1234-org/ino-platform-discovery/internal/request/github"
	"github.com/stretchr/testify/assert"
)

// TestNew provides coverage for the `New` function.
func TestNew(t *testing.T) {
	catalogEntry := New()
	assert.Equal(t, int(catalogcolumn.EnumEnd), len(catalogEntry))
}

// TestIsDuplicateNoClient provides basic coverage for the `IsDuplicate` function.
func TestIsDuplicateNoClient(t *testing.T) {
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
			testMessage:                  "It should return false when existing has invalid platform data.",
			incomingRepository:           "https://github.com/foo-owner/foo-repo",
			incomingRepositoryDataFolder: "foo-folder",
			incomingBranchName:           "foo-branch",
			existingRepository:           ":",
			existingRepositoryDataFolder: "foo-folder",
			existingBranchName:           "foo-branch",
			expected:                     false,
		},
		{
			testMessage:                  "It should return false when incoming is novel platform.",
			incomingRepository:           "https://github.com/bar-owner/bar-repo",
			incomingRepositoryDataFolder: "bar-folder",
			incomingBranchName:           "bar-branch",
			// Use non-resolvable host to avoid request client dependency.
			existingRepository:           "https://example.com/foo-owner/foo-repo",
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
			testMessage:                    "It should return false when existing has invalid index data.",
			incomingPackageIndexRepository: "https://github.com/foo-owner/foo-repo",
			incomingBoardsManagerURL:       "https://raw.githubusercontent.com/foo-owner/foo-repo/refs/heads/main/package_foo_index.json",
			incomingPackageIndexFolder:     "foo-folder",
			incomingPackageIndexBranch:     "foo-branch",
			existingPackageIndexRepository: ":",
			existingBoardsManagerURL:       "https://raw.githubusercontent.com/foo-owner/foo-repo/refs/heads/main/package_foo_index.json",
			existingPackageIndexFolder:     "foo-folder",
			existingPackageIndexBranch:     "foo-branch",
			expected:                       false,
		},
		{
			testMessage:                    "It should return false when incoming is novel index.",
			incomingPackageIndexRepository: "https://github.com/bar-owner/bar-repo",
			incomingBoardsManagerURL:       "https://raw.githubusercontent.com/bar-owner/bar-repo/refs/heads/main/package_bar_index.json",
			incomingPackageIndexFolder:     "bar-folder",
			incomingPackageIndexBranch:     "bar-branch",
			// Use non-resolvable host to avoid request client dependency.
			existingPackageIndexRepository: "https://example.com/foo-owner/foo-repo",
			existingBoardsManagerURL:       "https://example.com/foo-owner/foo-repo/refs/heads/main/package_foo_index.json",
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

// TestIsDuplicateWithClient provides coverage for the `IsDuplicate` function code with a request client dependency.
func TestIsDuplicateWithClient(t *testing.T) {
	if os.Getenv("GITHUB_TOKEN") == "" {
		t.Skip("Required GITHUB_TOKEN environment variable not defined.")
	}

	clients.Clients.GitHub = github.NewClient(os.Getenv("GITHUB_TOKEN"))

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
			testMessage:                  "It should return true when incoming is duplicate platform after normalization of existing.",
			incomingRepository:           "https://github.com/per1234-org/MouseTo",
			incomingRepositoryDataFolder: "foo-folder",
			incomingBranchName:           "foo-branch",
			existingRepository:           "https://github.com/per1234/MouseTo",
			existingRepositoryDataFolder: "foo-folder",
			existingBranchName:           "foo-branch",
			expected:                     true,
		},
		{
			testMessage:                    "It should return true when incoming is duplicate index after normalization of existing.",
			incomingPackageIndexRepository: "https://github.com/per1234-org/MouseTo",
			incomingBoardsManagerURL:       "https://raw.githubusercontent.com/foo-owner/foo-repo/refs/heads/main/package_foo_index.json",
			incomingPackageIndexFolder:     "foo-folder",
			incomingPackageIndexBranch:     "foo-branch",
			existingPackageIndexRepository: "https://github.com/per1234/MouseTo",
			existingBoardsManagerURL:       "https://raw.githubusercontent.com/foo-owner/foo-repo/refs/heads/main/package_foo_index.json",
			existingPackageIndexFolder:     "foo-folder",
			existingPackageIndexBranch:     "foo-branch",
			expected:                       true,
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

// Test_isDuplicateIndex provides coverage for the `isDuplicateIndex` function.
func Test_isDuplicateIndex(t *testing.T) {
	existing := New()
	existing[catalogcolumn.PackageIndexRepository] = "existing PackageIndexRepository"
	existing[catalogcolumn.BoardsManagerURL] = "https://raw.githubusercontent.com/foo-owner/foo-repo/refs/heads/main/package_foo_index.json"
	existing[catalogcolumn.PackageIndexFolder] = "existing PackageIndexFolder"
	existing[catalogcolumn.PackageIndexBranch] = "existing PackageIndexBranch"

	incoming := New()
	incoming[catalogcolumn.PackageIndexRepository] = "existing PackageIndexRepository"
	incoming[catalogcolumn.BoardsManagerURL] = "https://raw.githubusercontent.com/foo-owner/foo-repo/refs/heads/main/package_foo_index.json"
	incoming[catalogcolumn.PackageIndexFolder] = "existing PackageIndexFolder"
	incoming[catalogcolumn.PackageIndexBranch] = "existing PackageIndexBranch"

	assert.True(
		t,
		isDuplicateIndex(incoming, existing),
		"It should return true when incoming has duplicate index data.",
	)

	incoming[catalogcolumn.PackageIndexRepository] = "incoming PackageIndexRepository"

	assert.False(
		t,
		isDuplicateIndex(incoming, existing),
		"It should return false when incoming has a unique index repository.",
	)

	incoming[catalogcolumn.PackageIndexRepository] = "existing PackageIndexRepository"
	incoming[catalogcolumn.BoardsManagerURL] = "https://raw.githubusercontent.com/foo-owner/foo-repo/refs/heads/main/package_bar_index.json"
	assert.False(
		t,
		isDuplicateIndex(incoming, existing),
		"It should return false when incoming has a unique index filename.",
	)
}

// Test_isDuplicatePlatform provides coverage for the `isDuplicatePlatform` function.
func Test_isDuplicatePlatform(t *testing.T) {
	existing := New()
	existing[catalogcolumn.Repository] = "existing Repository"
	existing[catalogcolumn.RepositoryDataFolder] = "existing RepositoryDataFolder"
	existing[catalogcolumn.BranchName] = "existing BranchName"

	incoming := New()
	incoming[catalogcolumn.Repository] = "existing Repository"
	incoming[catalogcolumn.RepositoryDataFolder] = "existing RepositoryDataFolder"
	incoming[catalogcolumn.BranchName] = "existing BranchName"

	assert.True(
		t,
		isDuplicatePlatform(incoming, existing),
		"It should return true when incoming has duplicate platform data.",
	)

	incoming[catalogcolumn.Repository] = "incoming Repository"

	assert.False(
		t,
		isDuplicatePlatform(incoming, existing),
		"It should return false when incoming has unique platform data.",
	)
}

// Test_resolveURL provides coverage for the `resolveURL` function.
func Test_resolveURL(t *testing.T) {
	inputURL := ":"
	url, err := resolveURL(inputURL)
	assert.Equal(t, inputURL, url, "It should return the input URL when the URL is invalid.")
	assert.Error(t, err, "It should return an error when the URL is invalid.")

	if os.Getenv("GITHUB_TOKEN") == "" {
		t.Skip("Required GITHUB_TOKEN environment variable not defined.")
	}

	clients.Clients.GitHub = github.NewClient(os.Getenv("GITHUB_TOKEN"))

	inputURL = "https://github.com/per1234-org/nonexistent"
	url, err = resolveURL(inputURL)
	assert.Equal(t, inputURL, url, "It should return the input URL when provided a dead URL.")
	assert.Error(t, err, "It should return an error when provided a dead URL.")

	url, err = resolveURL("https://github.com/per1234/MouseTo")
	assert.Equal(t,
		"https://github.com/per1234-org/MouseTo",
		url,
		"It should return the resolved URL when passed a URL with redirects.",
	)
	assert.NoError(t, err, "It should return the resolved URL when passed a URL with redirects.")
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
