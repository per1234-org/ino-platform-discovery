// Package cli defines the command line interface.
package cli

import (
	"github.com/per1234-org/ino-platform-discovery/internal/command/root"
	"github.com/spf13/cobra"
)

// Root creates a new root CLI command
func Root() *cobra.Command {
	var command = &cobra.Command{
		Run:   root.Run,
		Short: "A tool to discover Arduino boards platforms",
		Use:   "ino-platform-discovery [--catalog <catalog path>] [--output <output file path>]",
	}

	command.PersistentFlags().String("catalog", "ino-hardware-package-list.tsv", "Path to the inoplatforms catalog file.")
	command.PersistentFlags().String("exclusions", "exclusions.yml", "Path to the exclusions file.")
	command.PersistentFlags().String("output", "discoveries.tsv", "Path of the discovery results file that should be written by the tool.")

	return command
}
