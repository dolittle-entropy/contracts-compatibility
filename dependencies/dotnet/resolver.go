package dotnet

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-semver/semver"
)

// DepsResolver represents a system that can resolve the version of a dependency from a '.deps.json' file
type DepsResolver struct {
	packageName string
}

// NewDepsResolverFor creates a new DepsResolver for a given package
func NewDepsResolverFor(name string) *DepsResolver {
	return &DepsResolver{
		packageName: name,
	}
}

// ResolveDependencyFromContents resolves the dependency version from the given '.deps.json' file contents
func (resolver *DepsResolver) ResolveDependencyFromContents(contents []byte) (*semver.Version, error) {
	info := dotnetDependencies{}
	if err := json.Unmarshal(contents, &info); err != nil {
		return nil, fmt.Errorf("could not parse .deps.json file: %w", err)
	}

	for _, assembly := range info.Targets {
		for _, target := range assembly {
			for dependency, dependencyVersion := range target.Dependencies {
				if dependency == resolver.packageName {
					parsed, err := semver.NewVersion(dependencyVersion)
					if err != nil {
						return nil, fmt.Errorf("could not parse .deps.json dependency version: %w", err)
					}

					return parsed, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no dependency found for %v", resolver.packageName)
}
