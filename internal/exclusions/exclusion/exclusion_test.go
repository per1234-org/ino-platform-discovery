package exclusion

import (
	"regexp"
	"testing"

	"github.com/per1234-org/ino-platform-discovery/internal/results/result"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result/content"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result/host"
	"github.com/stretchr/testify/assert"
)

// TestTypeMatch provides coverage for the `(*Type) Match` method.
func TestTypeMatch(t *testing.T) {
	testTables := []struct {
		testName  string
		exclusion Type
		result    result.Type
		expected  bool
	}{
		{
			testName: "Return true when exclusion matches index result.",
			exclusion: Type{
				Host:  regexp.MustCompile("^github.com$"),
				Name:  regexp.MustCompile("^fooname$"),
				Owner: regexp.MustCompile("^fooowner$"),
				Path:  regexp.MustCompile("package_foo_index.json"),
			},
			result: result.Type{
				Content:        content.Index,
				Host:           host.GitHub,
				RepositoryName: "fooname",
				Owner:          "fooowner",
				Path:           "foopath/package_foo_index.json",
			},
			expected: true,
		},
		{
			testName: "Return true when exclusion matches platform result.",
			exclusion: Type{
				Host:  regexp.MustCompile("^github.com$"),
				Name:  regexp.MustCompile("^fooname$"),
				Owner: regexp.MustCompile("^fooowner$"),
				Path:  regexp.MustCompile("^foopath$"),
			},
			result: result.Type{
				Content:        content.Platform,
				Host:           host.GitHub,
				RepositoryName: "fooname",
				Owner:          "fooowner",
				Path:           "foopath/boards.txt",
			},
			expected: true,
		},
		{
			testName: "Return true on exclusion mismatch.",
			exclusion: Type{
				Host:  regexp.MustCompile("^github.com$"),
				Name:  regexp.MustCompile("^fooname$"),
				Owner: regexp.MustCompile("^fooowner$"),
				Path:  regexp.MustCompile("^foopath$"),
			},
			result: result.Type{
				Content:        content.Index,
				Host:           host.GitHub,
				RepositoryName: "barname",
				Owner:          "barowner",
				Path:           "barpath/package_bar_index.json",
			},
			expected: false,
		},
	}

	for _, testTable := range testTables {
		assert.Equal(t, testTable.expected, testTable.exclusion.Match(testTable.result), testTable.testName)
	}
}
