// Package ghsearch contains code for making requests to the `/search/code` endpoint of the GitHub API.
package ghsearch

import (
	"context"
	"fmt"
	"net/url"
	"slices"

	gogithub "github.com/google/go-github/v79/github"
	"github.com/per1234-org/ino-platform-discovery/internal/request/github"
	"github.com/per1234-org/ino-platform-discovery/internal/results"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result/content"
	"github.com/sirupsen/logrus"
)

// Search searches for GitHub repositories that contain a package index and/or platform.
func Search(clientContext context.Context, client *gogithub.Client) (results.Type, error) {
	fmt.Println("Searching GitHub for package indexes...")
	results, err := indexes(clientContext, client)
	if err != nil {
		return nil, err
	}

	fmt.Println("Searching GitHub for platforms...")
	platformResults, err := platforms(clientContext, client)
	if err != nil {
		return nil, err
	}

	results.Merge(platformResults)

	return results, nil
}

// indexes searches for package indexes.
func indexes(clientContext context.Context, client *gogithub.Client) (results.Type, error) {
	results := results.Type{}

	var filenameSearchResponse *gogithub.CodeSearchResult
	// Unlike the website search, the API does not support filename patterns. This crude approach is all that is available
	// to find files with a `package_*_index.json` name pattern.
	query := "fork:true+in:path+language:json+package_+_index.json"
	doRequest := true
	for doRequest {
		logrus.Tracef("Making GitHub API /search/code endpoint request for %s", query)
		doRequest = false
		var err error
		filenameSearchResponse, _, err = client.Search.Code(
			clientContext,
			url.QueryEscape(query),
			nil,
		)

		if err != nil {
			if github.ShouldRetry(err) {
				// Retry request.
				doRequest = true
				continue
			}

			// Error is not recoverable.
			return results, err
		}
	}

	var contentSearchResponse *gogithub.CodeSearchResult
	// The path search is not precise, so likely to yield false positives. These should be eliminated by also checking for
	// the presence of distinctive key names in the file content.
	// We might expect to be able to accomplish this in a single query like:
	// `fork:true+in:file,path+language:json+package_+_index.json+archiveFileName+architecture+version+url`
	// However, that will only return results where all keywords are present in the file content. For this reason, it is
	// necessary to run two separate searches to accomplish this.
	query = "fork:true+language:json+archiveFileName+architecture+version+url"
	doRequest = true
	for doRequest {
		logrus.Tracef("Making GitHub API /search/code endpoint request for %s", query)
		doRequest = false
		var err error
		contentSearchResponse, _, err = client.Search.Code(
			clientContext,
			url.QueryEscape(query),
			nil,
		)

		if err != nil {
			if github.ShouldRetry(err) {
				// Retry request.
				doRequest = true
				continue
			}

			// Error is not recoverable.
			return results, err
		}
	}

	// Remove any results that were not returned by both searches
	searchResults := intersecting(filenameSearchResponse, contentSearchResponse)

	for _, searchResult := range searchResults {
		result := result.Type{
			Content: content.Type{
				Index: true,
			},
			IndexPath:      *searchResult.Path,
			Owner:          *searchResult.Repository.Owner.Login,
			RepositoryName: *searchResult.Repository.Name,
			RepositoryURL:  *searchResult.Repository.HTMLURL,
		}

		results = append(results, result)
	}

	return results, nil
}

// intersecting returns the list of search results that are present in both of the provided lists.
func intersecting(resultsA *gogithub.CodeSearchResult, resultsB *gogithub.CodeSearchResult) []*gogithub.CodeResult {
	return slices.DeleteFunc(
		resultsA.CodeResults,
		func(resultA *gogithub.CodeResult) bool {
			for _, resultB := range resultsB.CodeResults {
				if *resultB.Repository.HTMLURL == *resultA.Repository.HTMLURL {
					return false
				}
			}

			// This result is not present in both sets.
			return true
		},
	)
}

// platforms searches for platforms.
func platforms(clientContext context.Context, client *gogithub.Client) (results.Type, error) {
	results := results.Type{}

	var searchResponse *gogithub.CodeSearchResult
	query := "filename:boards.txt+fork:true+\".upload.tool\""
	doRequest := true
	for doRequest {
		logrus.Tracef("Making GitHub API /search/code endpoint request for %s", query)
		doRequest = false
		var err error
		searchResponse, _, err = client.Search.Code(
			clientContext,
			url.QueryEscape(query),
			nil,
		)

		if err != nil {
			if github.ShouldRetry(err) {
				// Retry request.
				doRequest = true
				continue
			}

			// Error is not recoverable.
			return results, err
		}
	}

	for _, searchResult := range searchResponse.CodeResults {
		result := result.Type{
			Content: content.Type{
				Platform: true,
			},
			PlatformFilePath: *searchResult.Path,
			RepositoryURL:    *searchResult.Repository.HTMLURL,
		}

		results = append(results, result)
	}

	return results, nil
}
