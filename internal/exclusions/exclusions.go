// Package exclusions contains code for excluding discoveries.
package exclusions

import (
	"github.com/arduino/go-paths-helper"
	"github.com/per1234-org/ino-platform-discovery/internal/exclusions/exclusion"
	"github.com/per1234-org/ino-platform-discovery/internal/exclusions/rawexclusion"
	"go.yaml.in/yaml/v4"
)

// Type is the type of the exclusion rules.
type Type []exclusion.Type

// Load returns a Type object populated with the data from the exclusions file.
func Load(path *paths.Path) (Type, error) {
	exclusions := Type{}

	fileContent, err := path.ReadFile()
	if err != nil {
		return exclusions, err
	}

	rawExclusions := []rawexclusion.Type{}
	err = yaml.Unmarshal(fileContent, &rawExclusions)
	if err != nil {
		return exclusions, err
	}

	for _, rawExclusion := range rawExclusions {
		exclusion, err := rawExclusion.ToExclusion()
		if err != nil {
			return exclusions, err
		}

		exclusions = append(exclusions, exclusion)
	}

	return exclusions, nil
}
