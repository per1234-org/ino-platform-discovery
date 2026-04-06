// Package results contains code related to the list of discoveries.
package results

import (
	"testing"

	"github.com/per1234-org/ino-platform-discovery/internal/catalog"
	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogcolumn"
	"github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogentry"
	"github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result/content"
	"github.com/stretchr/testify/assert"
)

// TestTypeDeduplicate provides coverage for the `(*Type) Deduplicate` method.
func TestTypeDeduplicate(t *testing.T) {
	results := Type{
		{
			RepositoryURL: "https://example.com/foo",
		},
		{
			RepositoryURL: "https://example.com/bar",
		},
		{
			RepositoryURL: "https://example.com/baz",
		},
	}

	fooCatalogEntry := catalogentry.New()
	fooCatalogEntry[catalogcolumn.Repository] = "https://example.com/foo"

	barCatalogEntry := catalogentry.New()
	barCatalogEntry[catalogcolumn.PackageIndexRepository] = "https://example.com/bar"

	quxCatalogEntry := catalogentry.New()
	quxCatalogEntry[catalogcolumn.PackageIndexRepository] = "https://example.com/qux"

	catalog := catalog.Type{
		fooCatalogEntry,
		barCatalogEntry,
		quxCatalogEntry,
	}

	assertion := Type{
		{
			RepositoryURL: "https://example.com/baz",
		},
	}

	results.Deduplicate(catalog)

	assert.Equal(t, assertion, results)
}

// TestTypeMerge provides coverage for the `(*Type) Merge` method.
func TestTypeMerge(t *testing.T) {
	targets := Type{
		{
			RepositoryURL: "https://example.com/foo",
		},
		{
			Content: content.Type{
				Index: true,
			},
			DefaultBranch: "barbranch",
			IndexPath:     "barindexpath",
			Owner:         "barownner",
			RepositoryURL: "https://example.com/bar",
		},
	}
	additions := Type{
		{
			Content: content.Type{
				Platform: true,
			},
			IndexPath:        "barindexpath",
			PlatformFilePath: "barplatformpath",
			RepositoryURL:    "https://example.com/bar",
		},
		{
			RepositoryURL: "https://example.com/baz",
		},
	}

	assertion := Type{
		{
			RepositoryURL: "https://example.com/foo",
		},
		{
			Content: content.Type{
				Index:    true,
				Platform: true,
			},
			DefaultBranch:    "barbranch",
			IndexPath:        "barindexpath",
			Owner:            "barownner",
			PlatformFilePath: "barplatformpath",
			RepositoryURL:    "https://example.com/bar",
		},
		{
			RepositoryURL: "https://example.com/baz",
		},
	}

	targets.Merge(additions)

	assert.Equal(t, assertion, targets)
}

// TestTypeFilter provides coverage for the `(*Type) Filter` method.
func TestTypeFilter(t *testing.T) {
	results := Type{
		{
			RepositoryData: &repo.Type{
				Fork:  true,
				Ahead: true,
			},
			RepositoryURL: "https://example.com/foo",
		},
		{
			RepositoryData: &repo.Type{
				Fork:  true,
				Ahead: false,
			},
			RepositoryURL: "https://example.com/bar",
		},
		{
			RepositoryData: &repo.Type{
				Fork:  false,
				Ahead: false,
			},
			RepositoryURL: "https://example.com/baz",
		},
	}

	assertion := Type{
		{
			RepositoryData: &repo.Type{
				Fork:  true,
				Ahead: true,
			},
			RepositoryURL: "https://example.com/foo",
		},
		{
			RepositoryData: &repo.Type{
				Fork:  false,
				Ahead: false,
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
			Content: content.Type{
				Index:    true,
				Platform: false,
			},
			DefaultBranch:    "foo-branch",
			IndexPath:        "foo-index-path/package_foo_index.json",
			Owner:            "foo-owner",
			PlatformFilePath: "foo-platform-path/boards.txt",
			RepositoryName:   "foo-repo",
			RepositoryURL:    "https://github.com/foo-owner/foo-repo",
		},
		{
			Content: content.Type{
				Index:    false,
				Platform: true,
			},
			DefaultBranch:    "bar-branch",
			IndexPath:        "bar-index-path/package_bar_index.json",
			Owner:            "bar-owner",
			PlatformFilePath: "bar-platform-path/boards.txt",
			RepositoryName:   "bar-repo",
			RepositoryURL:    "https://github.com/bar-owner/bar-repo",
		},
		{
			Content: content.Type{
				Index:    true,
				Platform: true,
			},
			DefaultBranch:    "baz-branch",
			IndexPath:        "baz-index-path/package_baz_index.json",
			Owner:            "baz-owner",
			PlatformFilePath: "baz-platform-path/boards.txt",
			RepositoryName:   "baz-repo",
			RepositoryURL:    "https://github.com/bar-owner/baz-repo",
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
		{
			"",
			"",
			"",
			"https://github.com/bar-owner/baz-repo",
			"https://raw.githubusercontent.com/baz-owner/baz-repo/refs/heads/baz-branch/baz-index-path/package_baz_index.json",
			"/baz-platform-path/",
			"baz-branch",
			"https://github.com/bar-owner/baz-repo",
			"/baz-index-path/",
			"baz-branch",
			"",
			"",
			"",
		},
	}

	catalog := results.ToCatalog()

	assert.Equal(t, assertion, catalog)
}

// TestType_find provides coverage for the `(*Type) find` method.
func TestType_find(t *testing.T) {
	results := Type{
		{
			RepositoryData: &repo.Type{
				Fork:  true,
				Ahead: true,
			},
			RepositoryURL: "https://example.com/foo",
		},
		{
			RepositoryData: &repo.Type{
				Fork:  true,
				Ahead: false,
			},
			RepositoryURL: "https://example.com/bar",
		},
		{
			RepositoryData: &repo.Type{
				Fork:  false,
				Ahead: false,
			},
			RepositoryURL: "https://example.com/baz",
		},
	}

	query := result.Type{
		RepositoryURL: "https://example.com/bar",
	}

	assert.Equal(t, 1, results.find(query), "It should return the index when query is present in results.")

	query.RepositoryURL = "https://example.com/qux"

	assert.Equal(t, -1, results.find(query), "It should return -1 when query is not present in results.")
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
