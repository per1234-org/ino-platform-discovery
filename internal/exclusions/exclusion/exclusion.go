// Package exclusion contains code related to individual exclusion rules.
package exclusion

import (
	"regexp"
	"strings"

	"github.com/per1234-org/ino-platform-discovery/internal/results/result"
	"github.com/per1234-org/ino-platform-discovery/internal/results/result/content"
)

// Type is the type of an exclusion rule.
type Type struct {
	// Host is the regular expression for the Git host.
	Host *regexp.Regexp
	// Name is the regular expression for the repository name.
	Name *regexp.Regexp
	// Owner is the regular expression for the repository owner.
	Owner *regexp.Regexp
	// Path is the regular expression for the path of the discovery.
	Path *regexp.Regexp
}

// Match determines whether the given result matches the given exclusion rule.
func (exclusion Type) Match(result result.Type) bool {
	var resultPath string
	switch result.Content {
	case content.Index:
		resultPath = result.Path
	case content.Platform:
		pathComponents := strings.Split(result.Path, "/")
		pathParentComponents := pathComponents[:len(pathComponents)-1]
		resultPath = strings.Join(pathParentComponents, "/")
	default:
		panic("unhandled content type")
	}

	return exclusion.Host.MatchString(result.Host.String()) &&
		exclusion.Name.MatchString(result.RepositoryName) &&
		exclusion.Owner.MatchString(result.Owner) &&
		exclusion.Path.MatchString(resultPath)
}
