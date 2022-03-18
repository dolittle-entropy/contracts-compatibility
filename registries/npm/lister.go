package npm

import (
	"dolittle.io/contracts-compatibility/http"
	"dolittle.io/contracts-compatibility/versioning"
	"fmt"
	"github.com/coreos/go-semver/semver"
)

// ReleaseLister represents a system that can list released versions of an NPM package
type ReleaseLister struct {
	packageName string
}

// NewReleaseListerFor creates a new ReleaseLister for an NPM package with the given name
func NewReleaseListerFor(name string) *ReleaseLister {
	return &ReleaseLister{
		packageName: name,
	}
}

// ListReleasedVersions lists the released version of the NPM package
func (lister *ReleaseLister) ListReleasedVersions() (semver.Versions, error) {
	info := packageInfo{}
	if err := http.GetJSON("https://registry.npmjs.org/"+lister.packageName, &info); err != nil {
		return nil, fmt.Errorf("could not get NPM package information: %w", err)
	}

	versions := make([]string, 0)
	for _, versionInfo := range info.Versions {
		if versionInfo.Name != lister.packageName {
			continue
		}
		versions = append(versions, versionInfo.Version)
	}
	return versioning.NewReleaseVersions(versions), nil
}
