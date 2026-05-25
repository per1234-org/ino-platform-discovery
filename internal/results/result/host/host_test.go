package host

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTypeString provides coverage for the `(Type) String` method.
func TestTypeString(t *testing.T) {
	assert.Equal(t, "github.com", GitHub.String())
}
