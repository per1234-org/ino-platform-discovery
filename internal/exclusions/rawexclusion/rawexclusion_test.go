package rawexclusion

import (
	"regexp"
	"testing"

	exclusion "github.com/per1234-org/ino-platform-discovery/internal/exclusions/exclusion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTypeToExclusion provides coverage for the `(*Type) ToExclusion` method.
func TestTypeToExclusion(t *testing.T) {
	rawexclusion := Type{
		Host:  "^github.com$",
		Name:  "^fooname$",
		Owner: "^fooowner$",
		Path:  "package_foo_index.json",
	}

	expected := exclusion.Type{
		Host:  regexp.MustCompile("^github.com$"),
		Name:  regexp.MustCompile("^fooname$"),
		Owner: regexp.MustCompile("^fooowner$"),
		Path:  regexp.MustCompile("package_foo_index.json"),
	}

	exclusion, err := rawexclusion.ToExclusion()
	require.NoError(t, err)

	assert.Equal(t, expected, exclusion)
}
