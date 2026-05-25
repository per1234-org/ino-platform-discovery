// Package host contains code related to the Git host of the result.
package host

// Type is the type of the Git host.
//
//go:generate go tool golang.org/x/tools/cmd/stringer -type=Type -linecomment
type Type int

// The line comments set the string representation for the host.
const (
	GitHub Type = iota // github.com
)
