package ownercache

import (
	"testing"

	"github.com/per1234-org/ino-platform-discovery/internal/results/repo"
	"github.com/stretchr/testify/assert"
)

// TestNew provides coverage for the `New` function.
func TestNew(t *testing.T) {
	assert.NotNil(t, New(), "It should return a non-nil object.")
}

// TestTypeGetSet provides coverage for the `(Type) Get` and `(Type) Set` methods.
func TestTypeGetSet(t *testing.T) {
	cache := New()
	_, cached := cache.Get("foo-repo")
	assert.False(t, cached, "Second return value should be false when the specified item is not in the cache.")

	setData := repo.Type{
		Ahead:         false,
		DefaultBranch: "main",
		Error:         nil,
		Fork:          true,
	}
	cache.Set("foo-repo", setData)

	getData, cached := cache.Get("foo-repo")
	assert.Equal(t, setData, getData, "It should return the stored repository data.")
	assert.True(t, cached, "Second return value should be true when the specified item is in the cache.")
}
