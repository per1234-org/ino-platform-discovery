// Package catalogentry contains code related to the individual catalog entries.
package catalogentry

import "github.com/per1234-org/ino-platform-discovery/internal/catalog/catalogcolumn"

// New returns a new catalog entry object.
func New() []string {
	return make([]string, catalogcolumn.EnumEnd)
}
