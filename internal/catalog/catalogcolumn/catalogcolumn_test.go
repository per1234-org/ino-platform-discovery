package catalogcolumn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTypeString provides coverage for the `(Type) String` method.
func TestTypeString(t *testing.T) {
	assert.Equal(t, "Boards Manager URL", BoardsManagerURL.String())
}
