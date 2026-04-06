package result

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTypeMerge provides coverage for the `(*Type) Merge` method.
func TestTypeMerge(t *testing.T) {
	target := Type{
		DefaultBranch:    "foo branch",
		IndexPath:        "",
		PlatformFilePath: "/foo-platform-folder/",
		RepositoryURL:    "",
	}

	addition := Type{
		DefaultBranch:    "foo branch",
		IndexPath:        "/foo-index-folder/",
		PlatformFilePath: "",
		RepositoryURL:    "https://example.com/foo-repo-url",
	}

	assertion := Type{
		DefaultBranch:    "foo branch",
		IndexPath:        "/foo-index-folder/",
		PlatformFilePath: "/foo-platform-folder/",
		RepositoryURL:    "https://example.com/foo-repo-url",
	}

	target.Merge(addition)
	assert.Equal(t, assertion, target)
}

// TestTypeSame provides coverage for the `(*Type) Same` method.
func TestTypeSame(t *testing.T) {
	resultA := Type{
		DefaultBranch:    "foo branch",
		IndexPath:        "",
		PlatformFilePath: "/foo-platform-folder/",
		RepositoryURL:    "https://example.com/foo-repo-url",
	}

	resultB := Type{
		DefaultBranch:    "foo branch",
		IndexPath:        "/foo-index-folder/",
		PlatformFilePath: "/foo-platform-folder/",
		RepositoryURL:    "https://example.com/foo-repo-url",
	}

	assert.True(t, resultA.Same(resultB))

	resultB.RepositoryURL = "https://example.com/bar-repo-url"

	assert.False(t, resultA.Same(resultB))
}
