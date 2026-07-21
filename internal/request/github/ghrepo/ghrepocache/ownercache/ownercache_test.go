package ownercache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew provides coverage for the `New` function.
func TestNew(t *testing.T) {
	assert.NotNil(t, New(), "It should return a non-nil object.")
}
