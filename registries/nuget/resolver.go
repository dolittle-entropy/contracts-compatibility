package nuget

import (
	"dolittle.io/contracts-compatibility/http"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"regexp"
	"strings"
)

// DependencyResolver represents a system that can resolve the version of a dependency for a specific version of an NuGet package
type DependencyResolver struct {
	packageName           string
	dependencyPackageName string
}

// NewDependencyResolverFor creates a new DependencyResolver for a given NuGet package and dependency
func NewDependencyResolverFor(name, dependencyName string) *DependencyResolver {
	return &DependencyResolver{
		packageName:           name,
		dependencyPackageName: dependencyName,
	}
}

// ResolveDependencyForVersion resolves the dependency version for the given package version
func (resolver *DependencyResolver) ResolveDependencyForVersion(version *semver.Version) (*semver.Version, error) {
	packageLinkInfo := catalogPackageLink{}
	if err := http.GetJSON("https://api.nuget.org/v3/registration5-semver1/"+strings.ToLower(resolver.packageName)+"/"+strings.ToLower(version.String())+".json", &packageLinkInfo); err != nil {
		return nil, fmt.Errorf("could not get NuGet package link information: %w", err)
	}

	info := catalogPackageEntry{}
	if err := http.GetJSON(packageLinkInfo.Entry, &info); err != nil {
		return nil, fmt.Errorf("could not get NuGet package entry information: %w", err)
	}

	for _, dependencyGroup := range info.DependencyGroups {
		for _, dependency := range dependencyGroup.Dependencies {
			if dependency.Id != resolver.dependencyPackageName {
				continue
			}

			matches := nugetDependencyRangeExpression.FindStringSubmatch(dependency.DependencyRange)
			if len(matches) != 2 {
				return nil, fmt.Errorf("could not parse dependency range %v", dependency.DependencyRange)
			}

			contractsVersion, err := semver.NewVersion(matches[1])
			if err != nil {
				return nil, fmt.Errorf("could not parse NuGet dependency version: %w", err)
			}

			return contractsVersion, nil
		}
	}

	return nil, fmt.Errorf("no dependency found for %v", resolver.dependencyPackageName)
}

var nugetDependencyRangeExpression = regexp.MustCompile("^\\[([^,]+), \\)$")
