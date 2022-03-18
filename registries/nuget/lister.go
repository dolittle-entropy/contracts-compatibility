package nuget

import (
	"dolittle.io/contracts-compatibility/http"
	"dolittle.io/contracts-compatibility/versioning"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"strings"
)

// ReleaseLister represents a system that can list released versions of an NuGet package
type ReleaseLister struct {
	packageName string
}

// NewReleaseListerFor creates a new ReleaseLister for an NuGet package with the given name
func NewReleaseListerFor(name string) *ReleaseLister {
	return &ReleaseLister{
		packageName: name,
	}
}

// ListReleasedVersions lists the released version of the NuGet package
func (lister *ReleaseLister) ListReleasedVersions() (semver.Versions, error) {
	registrationInfo := catalogRegistration{}
	if err := http.GetJSON("https://api.nuget.org/v3/registration5-semver1/"+strings.ToLower(lister.packageName)+"/index.json", &registrationInfo); err != nil {
		return nil, fmt.Errorf("could not get NuGet registration information: %w", err)
	}

	versions := make([]string, 0)

	for _, pageInfo := range registrationInfo.Pages {
		for _, packageInfo := range pageInfo.Packages {
			if packageInfo.Entry.Id != lister.packageName {
				continue
			}
			versions = append(versions, packageInfo.Entry.Version)
		}
	}

	return versioning.NewReleaseVersions(versions), nil
}
