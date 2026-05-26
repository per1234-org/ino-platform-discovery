// Package rawexclusion contains code related to the raw exclusion rule data.
package rawexclusion

import (
	"regexp"

	exclusion "github.com/per1234-org/ino-platform-discovery/internal/exclusions/exclusion"
)

// Type is the type of the raw data for an exclusion rule.
type Type struct {
	// Host is the Git host.
	Host string
	// Name is the repository name.
	Name string
	// Owner is the repository owner.
	Owner string
	// Path is the path of the discovery.
	Path string
}

// ToExclusion converts the raw exclusion rule data to an exclusion rule.
func (rawexclusion Type) ToExclusion() (exclusion.Type, error) {
	exclusion := exclusion.Type{}

	hostRegexp, err := regexp.Compile(rawexclusion.Host)
	if err != nil {
		return exclusion, err
	}
	exclusion.Host = hostRegexp

	nameRegexp, err := regexp.Compile(rawexclusion.Name)
	if err != nil {
		return exclusion, err
	}
	exclusion.Name = nameRegexp

	ownerRegexp, err := regexp.Compile(rawexclusion.Owner)
	if err != nil {
		return exclusion, err
	}
	exclusion.Owner = ownerRegexp

	pathRegexp, err := regexp.Compile(rawexclusion.Path)
	if err != nil {
		return exclusion, err
	}
	exclusion.Path = pathRegexp

	return exclusion, nil
}

// UnmarshalYAML implements a custom unmarshal function that sets default values, for use by
// go.yaml.in/yaml/v4.Unmarshal.
func (rawexclusion *Type) UnmarshalYAML(unmarshal func(any) error) error {
	type typeType Type
	withDefaults := typeType{
		Name: ".*",
		Path: ".*",
	}
	if err := unmarshal(&withDefaults); err != nil {
		return err
	}

	*rawexclusion = Type(withDefaults)

	return nil
}
