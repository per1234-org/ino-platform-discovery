// Package catalogcolumn contains code related to the columns of the catalog spreadsheet
package catalogcolumn

// Type is the type of the columns.
//
//go:generate go tool golang.org/x/tools/cmd/stringer -type=Type -linecomment
type Type int

/*
This defines the order of the columns.

The line comments set the heading string for each column.
*/
const (
	Name                   Type = iota // Name
	Vendor                             // Vendor
	Architecture                       // Architecture
	Repository                         // Repository
	BoardsManagerURL                   // Boards Manager URL
	RepositoryDataFolder               // Repository Data Folder
	BranchName                         // Branch Name
	PackageIndexRepository             // Package Index Repository
	PackageIndexFolder                 // Package Index Folder
	PackageIndexBranch                 // Package Index Branch
	Reference                          // Reference
	Notes                              // Notes
	Suppress                           // Suppress
	EnumEnd
)
