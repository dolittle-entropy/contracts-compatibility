package npm

import (
	"dolittle.io/contracts-compatibility/http"
	"fmt"
	"github.com/coreos/go-semver/semver"
)

// DependencyResolver represents a system that can resolve the version of a dependency for a specific version of an NPM package
type DependencyResolver struct {
	packageName           string
	dependencyPackageName string
}

// NewDependencyResolverFor creates a new DependencyResolver for a given NPM package and dependency
func NewDependencyResolverFor(name, dependencyName string) *DependencyResolver {
	return &DependencyResolver{
		packageName:           name,
		dependencyPackageName: dependencyName,
	}
}

// ResolveDependencyForVersion resolves the dependency version for the given package version
func (resolver *DependencyResolver) ResolveDependencyForVersion(version *semver.Version) (*semver.Version, error) {
	info := packageVersionInfo{}
	if err := http.GetJSON("https://registry.npmjs.org/"+resolver.packageName+"/"+version.String(), &info); err != nil {
		return nil, fmt.Errorf("could not get NPM package version information: %w", err)
	}

	for dependency, dependencyVersion := range info.Dependencies {
		if dependency != resolver.dependencyPackageName {
			continue
		}

		parsed, err := semver.NewVersion(dependencyVersion)
		if err != nil {
			return nil, fmt.Errorf("could not parse NPM dependency version: %w", err)
		}

		return parsed, nil
	}

	return nil, fmt.Errorf("no dependency found for %v", resolver.dependencyPackageName)
}
