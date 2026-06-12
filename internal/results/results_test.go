// Package results contains code related to the list of discoveries.
package results

import (
	"testing"

	"github.com/per1234-org/ino-platform-discovery/internal/catalog"
	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogcolumn"
	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogentry"
	"github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result/content"
	"github.com/stretchr/testify/assert"
)

// TestTypeDeduplicate provides coverage for the `(*Type) Deduplicate` method.
func TestTypeDeduplicate(t *testing.T) {
	results := Type{
		{
			Content: content.Platform,
			Path:    "foopath/boards.txt",
			RepositoryData: repo.Type{
				DefaultBranch: "foobranch",
			},
			RepositoryURL: "https://example.com/duplicate-platform",
		},
		{
			Content: content.Index,
			Path:    "foopath/package_foo_index.json",
			RepositoryData: repo.Type{
				DefaultBranch: "foobranch",
			},
			RepositoryURL: "https://example.com/duplicate-index",
		},
		{
			Content: content.Platform,
			Path:    "foopath/boards.txt",
			RepositoryData: repo.Type{
				DefaultBranch: "foobranch",
			},
			RepositoryURL: "https://example.com/novel-platform",
		},
	}

	fooCatalogEntry := catalogentry.New()
	fooCatalogEntry[catalogcolumn.Repository] = "https://example.com/duplicate-platform"
	fooCatalogEntry[catalogcolumn.RepositoryDataFolder] = "/foopath/"
	fooCatalogEntry[catalogcolumn.BranchName] = "foobranch"

	barCatalogEntry := catalogentry.New()
	barCatalogEntry[catalogcolumn.PackageIndexRepository] = "https://example.com/duplicate-index"
	barCatalogEntry[catalogcolumn.PackageIndexFolder] = "/foopath/"
	barCatalogEntry[catalogcolumn.PackageIndexBranch] = "foobranch"

	bazCatalogEntry := catalogentry.New()
	bazCatalogEntry[catalogcolumn.PackageIndexRepository] = "https://example.com/some-index"
	bazCatalogEntry[catalogcolumn.PackageIndexFolder] = "/foopath/"
	bazCatalogEntry[catalogcolumn.PackageIndexBranch] = "foobranch"

	catalog := catalog.Type{
		fooCatalogEntry,
		barCatalogEntry,
		bazCatalogEntry,
	}

	assertion := Type{
		{
			Content: content.Platform,
			Path:    "foopath/boards.txt",
			RepositoryData: repo.Type{
				DefaultBranch: "foobranch",
			},
			RepositoryURL: "https://example.com/novel-platform",
		},
	}

	results.Deduplicate(catalog)

	assert.Equal(t, assertion, results)
}

// TestTypeFilter provides coverage for the `(*Type) Filter` method.
func TestTypeFilter(t *testing.T) {
	unsupplementedResults := Type{
		{
			RepositoryURL: "https://example.com/foo",
		},
	}

	assert.Panics(
		t,
		func() {
			unsupplementedResults.Filter()
		},
		"Should panic if results have not been supplemented",
	)

	results := Type{
		{
			RepositoryData: repo.Type{
				DefaultBranch: "foo-branch",
				Fork:          true,
				Ahead:         true,
			},
			RepositoryURL: "https://example.com/foo",
		},
		{
			RepositoryData: repo.Type{
				DefaultBranch: "bar-branch",
				Fork:          true,
				Ahead:         false,
			},
			RepositoryURL: "https://example.com/bar",
		},
		{
			RepositoryData: repo.Type{
				DefaultBranch: "baz-branch",
				Fork:          false,
				Ahead:         false,
			},
			RepositoryURL: "https://example.com/baz",
		},
	}

	assertion := Type{
		{
			RepositoryData: repo.Type{
				DefaultBranch: "foo-branch",
				Fork:          true,
				Ahead:         true,
			},
			RepositoryURL: "https://example.com/foo",
		},
		{
			RepositoryData: repo.Type{
				DefaultBranch: "baz-branch",
				Fork:          false,
				Ahead:         false,
			},
			RepositoryURL: "https://example.com/baz",
		},
	}

	results.Filter()

	assert.Equal(t, assertion, results)
}

// TestTypeToCatalog provides coverage for the `(*Type) ToCatalog` method.
func TestTypeToCatalog(t *testing.T) {
	results := Type{
		{
			Content: content.Index,
			Owner:   "foo-owner",
			Path:    "foo-index-path/package_foo_index.json",
			RepositoryData: repo.Type{
				DefaultBranch: "foo-branch",
			},
			RepositoryName: "foo-repo",
			RepositoryURL:  "https://github.com/foo-owner/foo-repo",
		},
		{
			Content: content.Platform,
			Owner:   "bar-owner",
			Path:    "bar-platform-path/boards.txt",
			RepositoryData: repo.Type{
				DefaultBranch: "bar-branch",
			},
			RepositoryName: "bar-repo",
			RepositoryURL:  "https://github.com/bar-owner/bar-repo",
		},
	}

	assertion := catalog.Type{
		{
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
		},
		{
			"",
			"",
			"",
			"https://github.com/bar-owner/bar-repo",
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

	catalog := results.ToCatalog()

	assert.Equal(t, assertion, catalog)
}
