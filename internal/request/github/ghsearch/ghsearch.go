// Package ghsearch contains code for making requests to the `/search/code` endpoint of the GitHub API.
package ghsearch

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	gogithub "github.com/google/go-github/v79/github"
	"github.com/per1234-org/ino-platform-discovery/internal/data"
	"github.com/per1234-org/ino-platform-discovery/internal/request/github"
	"github.com/per1234-org/ino-platform-discovery/internal/results"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result/content"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result/host"
	"github.com/sirupsen/logrus"
)

// Search searches for GitHub repositories that contain a package index and/or platform.
func Search(client *gogithub.Client) (results.Type, error) {
	results, err := indexes(client)
	if err != nil {
		return nil, err
	}

	platformResults, err := platforms(client)
	if err != nil {
		return nil, err
	}

	/*
		Some index and platform results may be in the same repository. In cases where an index supports a platform, they
		should be merged into the same result. However, the fact they are both in the same repository is no guarantee of
		this association (because the repo may contain multiple platforms with separate indexes). Identification of
		associations between discovery results and merging will need to be performed by the human user.
	*/
	results = append(results, platformResults...)
	return results, nil
}

// indexes searches for package indexes.
func indexes(client *gogithub.Client) (results.Type, error) {
	results := results.Type{}
	var err error

	/*
		See: https://docs.github.com/search-github/searching-on-github/searching-code

		Unlike the website search, the API does not support filename patterns. This crude approach is all that is available
		to find files with a `package_*_index.json` name pattern.
	*/
	query := "in:path language:json package_ _index.json"
	fmt.Println("Searching GitHub for package indexes...")
	searchResults, err := search(client, query)
	if err != nil {
		return results, err
	}

	if len(searchResults) == 0 {
		// The raw search will always return results, so this is an "impossible" outcome that indicates a bug in the code.
		panic("no results from index search")
	}

	fmt.Println("Validating package index search results...")
	for _, searchResult := range searchResults {
		/*
			The code search query syntax doesn't provide any mechanism for specifying an exact filename format, so might
			return invalid results which can be identified by having a noncompliant filename.

			An equivalent check is done for platform results in `results.Type.Prefilter`
		*/
		if !verifyIndexFilename(*searchResult.Name) {
			// This search result is does not have a valid package index filename, so exclude it.
			continue
		}

		/*
			Non-index files might have a filename that happens to match the package index filename format. These should be
			eliminated by also checking for the presence of distinctive key names in the file content. We might expect to be
			able to accomplish this in a single search query like:
			`fork:true+in:file,path+language:json+package_+_index.json+archiveFileName+architecture+version+url`
			However, that will only return results where all keywords are present in the file content. For this reason, the
			approach taken is to instead request the file content for each search result, and check whether it looks like an
			index.
		*/
		var contentReader io.ReadCloser
		var response *gogithub.Response
		for {
			logrus.WithFields(
				logrus.Fields{
					"Repo": fmt.Sprintf(
						"%s/%s",
						*searchResult.Repository.Owner.Login,
						*searchResult.Repository.Name,
					),
					"Path": *searchResult.Path,
				},
			).Trace(
				"Making GitHub API /repos/{owner}/{repo}/contents/{path} endpoint request.",
			)
			var err error
			contentReader, response, err = client.Repositories.DownloadContents(
				context.Background(),
				*searchResult.Repository.Owner.Login,
				*searchResult.Repository.Name,
				*searchResult.Path,
				nil,
			)

			if err != nil {
				err := github.HandleRateLimiting(err)
				if err != nil {
					// Error is not recoverable.
					return results, err
				}

				// Retry request.
				continue
			}

			// It is possible for the download to result in a failed response when the returned error is nil.
			// See: https://pkg.go.dev/github.com/google/go-github/v79@v79.0.0/github#RepositoriesService.DownloadContents
			if response.StatusCode != http.StatusOK {
				// Retry request.
				continue
			}

			// Request was successful.
			defer func() {
				if closeError := contentReader.Close(); closeError != nil {
					err = errors.Join(err, closeError)
				}
			}()

			break
		}

		verified, err := verifyIndex(contentReader)
		if err != nil {
			return results, err
		}

		if !verified {
			// This search result is not an index, so exclude it.
			continue
		}

		result := result.Type{
			Content:        content.Index,
			Filename:       *searchResult.Name,
			Host:           host.GitHub,
			Owner:          *searchResult.Repository.Owner.Login,
			Path:           *searchResult.Path,
			RepositoryName: *searchResult.Repository.Name,
			RepositoryURL:  *searchResult.Repository.HTMLURL,
		}

		results = append(results, result)
	}

	if len(results) == 0 {
		// The raw search will always return verifiable indexes, so this is an "impossible" outcome that indicates a bug in
		// the code.
		panic("all index search results failed content verification")
	}

	return results, err
}

// platforms searches for platforms.
func platforms(client *gogithub.Client) (results.Type, error) {
	results := results.Type{}

	// See: https://docs.github.com/search-github/searching-on-github/searching-code
	query := fmt.Sprintf("filename:%s \".upload.tool\"", data.PlatformIndicatorFile)
	fmt.Println("Searching GitHub for platforms...")
	searchResults, err := search(client, query)
	if err != nil {
		return results, err
	}

	for _, searchResult := range searchResults {
		result := result.Type{
			Content:        content.Platform,
			Filename:       *searchResult.Name,
			Host:           host.GitHub,
			Path:           *searchResult.Path,
			Owner:          *searchResult.Repository.Owner.Login,
			RepositoryName: *searchResult.Repository.Name,
			RepositoryURL:  *searchResult.Repository.HTMLURL,
		}

		results = append(results, result)
	}

	if len(results) == 0 {
		// The raw search will always return results, so this is an "impossible" outcome that indicates a bug in the code.
		panic("no results from platform search")
	}

	return results, nil
}

// search makes requests to the `/search/code` endpoint of the GitHub API.
func search(client *gogithub.Client, query string) ([]*gogithub.CodeResult, error) {
	results := []*gogithub.CodeResult{}
	requestOptions := &gogithub.SearchOptions{
		ListOptions: gogithub.ListOptions{
			// Pages are 1-indexed.
			Page:    1,
			PerPage: 100,
		},
	}

	logrus.Tracef("Making GitHub API /search/code endpoint request for %s", query)
	for {
		logrus.Tracef("Requesting results page %v", requestOptions.Page)
		result, response, err := client.Search.Code(
			context.Background(),
			query,
			requestOptions,
		)

		if err != nil {
			err := github.HandleRateLimiting(err)
			if err != nil {
				// Error is not recoverable.
				return results, err
			}

			// Retry request.
			continue
		}

		// Request was successful.
		results = append(results, result.CodeResults...)

		// Handle pagination.
		if response.NextPage == 0 {
			// Pagination completed.
			break
		}
		requestOptions.Page = response.NextPage
	}

	return results, nil
}

// verifyIndex determines whether the given content is intended to be an Arduino package index.
func verifyIndex(reader io.Reader) (bool, error) {
	type criterion struct {
		content string
		found   bool
	}
	criteria := []criterion{
		{
			content: "\"architecture\"",
		},
		{
			content: "\"archiveFileName\"",
		},
		{
			content: "\"url\"",
		},
	}

	// A crude text search approach is used instead of proper JSON parsing because the goal is to discover everything
	// intended to be a package index, even in the case where the file does not have a valid JSON format.
	scanner := bufio.NewScanner(reader)

	// Scan the content, line by line.
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			if errors.Is(err, bufio.ErrTooLong) {
				// `bufio.Scanner` errors if the line length exceeds bufio.MaxScanTokenSize. This specific error should not be
				// considered a fault condition, but does indicate this is not a package index.
				return false, nil
			}

			return false, err
		}

		line := scanner.Text()
		for criterionIndex, criterion := range criteria {
			if strings.Contains(line, criterion.content) {
				criteria[criterionIndex].found = true
				break
			}
		}
	}

	for _, criterion := range criteria {
		if !criterion.found {
			return false, nil
		}
	}

	return true, nil
}

// verifyIndexFilename determines whether the filename of the result is that of a package index.
func verifyIndexFilename(filename string) bool {
	// See: https://arduino.github.io/arduino-cli/latest/package_index_json-specification/#naming-of-the-json-index-file
	return strings.HasPrefix(filename, "package_") && strings.HasSuffix(filename, "_index.json")
}
