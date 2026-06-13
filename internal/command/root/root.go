// Package root implements the root CLI command.
package root

import (
	"fmt"
	"os"

	"github.com/arduino/go-paths-helper"
	"github.com/per1234-org/ino-platform-discovery/internal/catalog"
	"github.com/per1234-org/ino-platform-discovery/internal/exclusions"
	"github.com/per1234-org/ino-platform-discovery/internal/feedback"
	"github.com/per1234-org/ino-platform-discovery/internal/request"
	"github.com/per1234-org/ino-platform-discovery/internal/request/github"
	"github.com/spf13/cobra"
)

// Run runs the root CLI command.
func Run(command *cobra.Command, _ []string) {
	if err := validateUserInput(command); err != nil {
		feedback.Error(err)
		os.Exit(1)
	}

	// Load the data from the catalog file
	catalogArg, err := command.Flags().GetString("catalog")
	if err != nil {
		panic(err)
	}
	catalog, err := catalog.Load(paths.New(catalogArg))
	if err != nil {
		feedback.Error(fmt.Errorf("while loading catalog file %s: %s", catalogArg, err))
		os.Exit(1)
	}

	// Load the data from the exclusions file.
	exclusionsArg, err := command.Flags().GetString("exclusions")
	if err != nil {
		panic(err)
	}
	exclusions, err := exclusions.Load(paths.New(exclusionsArg))
	if err != nil {
		feedback.Error(fmt.Errorf("while loading exclusions file %s: %s", exclusionsArg, err))
		os.Exit(1)
	}

	githubContext, githubClient := github.NewClient(os.Getenv("GITHUB_TOKEN"))

	// Search GitHub for repositories that appear to contain a platform or package index.
	searchResults, err := request.Search(githubContext, githubClient)
	if err != nil {
		feedback.Error(fmt.Errorf("while searching: %s", err))
		os.Exit(1)
	}

	// Remove results the search data indicates as invalid.
	searchResults.Prefilter()

	// Remove excluded results.
	searchResults.Exclude(exclusions)

	// Obtain additional data for each of the results.
	fmt.Println("Obtaining supplemental data for discoveries...")
	if err := request.Supplement(githubContext, githubClient, &searchResults); err != nil {
		feedback.Error(fmt.Errorf("while supplementing results: %s", err))
		os.Exit(1)
	}

	// Remove results the supplemental data indicates as invalid.
	searchResults.FilterSupplemented()

	/*
		Remove results already present in the catalog.
		This must be performed after supplementing the results, as the branch name is one of the items compared to determine
		whether a result is a duplicate, and that is added by supplementation.
	*/
	searchResults.Deduplicate(catalog)

	if len(searchResults) == 0 {
		feedback.Warning(fmt.Errorf("no discoveries were made"))
		// Not treated as a failure because it may occur under normal operating conditions: If all discoveries from a
		// previous run were added to the catalog after a previous run, and no new projects have been created since.
		os.Exit(0)
	}

	// Generate output file.
	discoveries := searchResults.ToCatalog()

	outputArg, err := command.Flags().GetString("output")
	if err != nil {
		panic(err)
	}
	if err := discoveries.Write(paths.New(outputArg)); err != nil {
		feedback.Error(fmt.Errorf("while writing discoveries output file %s: %s", outputArg, err))
		os.Exit(1)
	}

	fmt.Printf("Discovery finished successfully. Results saved to: %s\n", outputArg)
}

// validateUserInput validates the user input for the command.
func validateUserInput(command *cobra.Command) error {
	if os.Getenv("GITHUB_TOKEN") == "" {
		return fmt.Errorf("environment variable GITHUB_TOKEN not set")
	}

	catalogArg, err := command.Flags().GetString("catalog")
	if err != nil {
		panic(err)
	}

	if exist, err := paths.New(catalogArg).ExistCheck(); !exist {
		if err == nil {
			return fmt.Errorf("file not found: %s", catalogArg)
		}

		return fmt.Errorf("unable to access %s: %s", catalogArg, err)
	}

	exclusionsArg, err := command.Flags().GetString("exclusions")
	if err != nil {
		panic(err)
	}

	if exist, err := paths.New(exclusionsArg).ExistCheck(); !exist {
		if err == nil {
			return fmt.Errorf("file not found: %s", exclusionsArg)
		}

		return fmt.Errorf("unable to access %s: %s", exclusionsArg, err)
	}

	outputArg, err := command.Flags().GetString("output")
	if err != nil {
		panic(err)
	}

	outputParentPath := paths.New(outputArg).Parent()
	if exist, err := outputParentPath.ExistCheck(); !exist {
		if err == nil {
			return fmt.Errorf("parent folder of output path not found: %s", outputParentPath)
		}

		return fmt.Errorf("unable to access parent folder of output path %s: %s", outputParentPath, err)
	}

	return nil
}
