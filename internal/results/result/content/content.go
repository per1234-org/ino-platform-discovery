// Package content contains code related to the content of the result.
package content

// Type is the type of the discovered content type.
type Type int

const (
	// Index is a package index.
	Index Type = iota
	// Platform is a boards platform.
	Platform
)
