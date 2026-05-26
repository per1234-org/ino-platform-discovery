package exclusions

import (
	"regexp"
	"testing"

	"github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoad provides coverage for the `Load` function.
func TestLoad(t *testing.T) {
	expected := Type{
		{
			Host:  regexp.MustCompile("foohost"),
			Name:  regexp.MustCompile("fooname"),
			Owner: regexp.MustCompile("fooowner"),
			Path:  regexp.MustCompile("foopath"),
		},
		{
			Host:  regexp.MustCompile("barhost"),
			Name:  regexp.MustCompile(".*"),
			Owner: regexp.MustCompile("barowner"),
			Path:  regexp.MustCompile(".*"),
		},
	}

	workingDirectory, err := paths.Getwd()
	require.NoError(t, err)
	exclusionsPath := workingDirectory.Join("testdata", "TestLoad", "exclusions.yml")
	exclusions, err := Load(exclusionsPath)
	require.NoError(t, err)

	assert.Equal(t, expected, exclusions)
}
