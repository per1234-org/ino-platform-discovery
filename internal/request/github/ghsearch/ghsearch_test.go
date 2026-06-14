package ghsearch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_verifyIndex provides coverage for the `verifyIndex` function.
func Test_verifyIndex(t *testing.T) {
	reader := strings.NewReader(`
{
  "packages": [
    {
      "name": "foo",
      "maintainer": "Foo",
      "websiteURL": "https://example.com/",
      "email": "nope@example.com",
      "help": {
        "online": "https://example.com"
      },
      "platforms": [
        {
          "name": "Bar Boards",
          "architecture": "bar",
          "version": "1.2.3",
          "category": "Contributed",
          "url": "https://example.com/bar-1.2.3.tar.bz2",
          "archiveFileName": "bar-1.2.3.tar.bz2",
          "checksum": "SHA-256:6ae7000b1b6c004c4a208d6d924a8046b417580e4a10cfd3f88930b6660a51db",
          "size": "7161883",
          "help": {
            "online": "https://example.com"
          },
          "boards": [
            {
              "name": "Classic Bar"
            },
            {
              "name": "Bar++"
            }
          ],
          "toolsDependencies": []
        }
      ]
    }
  ]
}
`,
	)

	verified, err := verifyIndex(reader)
	require.NoError(t, err)
	assert.True(t, verified)

	reader = strings.NewReader(`"foo": {}`)
	verified, err = verifyIndex(reader)
	require.NoError(t, err)
	assert.False(t, verified)
}

// Test_verifyIndexFilename provides coverage for the `verifyIndexFilename` function.
func Test_verifyIndexFilename(t *testing.T) {
	assert.True(t, verifyIndexFilename("package_foo_index.json"))

	assert.False(t, verifyIndexFilename("foo_package_bar_index.json"))
}
