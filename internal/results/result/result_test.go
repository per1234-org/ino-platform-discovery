package result

import (
	"testing"

	"github.com/per1234-org/ino-platform-discovery/internal/catalog"
	"github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result/content"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result/host"
	"github.com/stretchr/testify/assert"
)

// TestIsDuplicate provides coverage for the `(Type) IsDuplicate` method.
func TestIsDuplicate(t *testing.T) {
	theResult := Type{
		Content: content.Platform,
		Path:    "foo-platform-path/boards.txt",
		RepositoryData: repo.Type{
			DefaultBranch: "foo-branch",
		},
		RepositoryURL: "https://example.com/foo-owner/foo-repo",
	}

	theCatalog := catalog.Type{
		{
			"",
			"",
			"",
			"https://example.com/bar-owner/bar-repo",
			"",
			"/bar-platform-path/",
			"bar-branch",
			"",
			"",
			"",
			"",
			"",
			"",
		},
		{
			"",
			"",
			"",
			"https://example.com/foo-owner/foo-repo",
			"",
			"/foo-platform-path/",
			"foo-branch",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	}

	assert.True(t, theResult.IsDuplicate(theCatalog), "It should return true when matching catalog entry present")

	theCatalog = catalog.Type{
		{
			"",
			"",
			"",
			"https://example.com/bar-owner/bar-repo",
			"",
			"/bar-platform-path/",
			"bar-branch",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	}

	assert.False(t, theResult.IsDuplicate(theCatalog), "It should return false when no matching catalog entry present")
}

// TestTypeToCatalogEntry provides coverage for the `(*Type) ToCatalogEntry` method.
func TestTypeToCatalogEntry(t *testing.T) {
	result := Type{
		Content:        content.Index,
		Host:           host.GitHub,
		Owner:          "foo-owner",
		Path:           "foo-index-path/package_foo_index.json",
		RepositoryName: "foo-repo",
		RepositoryURL:  "https://github.com/foo-owner/foo-repo",
	}

	assert.Panics(
		t,
		func() {
			result.ToCatalogEntry()
		},
		"Should panic if result has not been supplemented",
	)

	result.RepositoryData = repo.Type{
		DefaultBranch: "foo-branch",
	}

	assertion := []string{
		"",
		"",
		"",
		"",
		"https://raw.githubusercontent.com/foo-owner/foo-repo/refs/heads/foo-branch/foo-index-path/package_foo_index.json",
		"",
		"",
		"https://github.com/foo-owner/foo-repo",
		"/foo-index-path/",
		"foo-branch",
		"",
		"",
		"",
	}

	assert.Equal(t, assertion, result.ToCatalogEntry())
}

// Test_toCatalogPath provides coverage for the `toCatalogPath` function.
func Test_toCatalogPath(t *testing.T) {
	testTables := []struct {
		testName   string
		resultPath string
		assertion  string
	}{
		{
			"Should return `/` when result path is in root.",
			"package_foo_index.json",
			"/",
		},
		{
			"Should return parent path w/ leading and trailing separator when result path is in subfolder.",
			"foo/platform/path/boards.txt",
			"/foo/platform/path/",
		},
	}

	for _, testTable := range testTables {
		assert.Equal(t, testTable.assertion, toCatalogPath(testTable.resultPath), testTable.testName)
	}
}
