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
