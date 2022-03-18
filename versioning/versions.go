package versioning

import (
	"github.com/coreos/go-semver/semver"
)

// NewReleaseVersions parses the version strings and returns valid non-prerelease values
func NewReleaseVersions(versions []string) semver.Versions {
	parsed := make(semver.Versions, 0)

	for _, version := range versions {
		parsedVersion, err := NewReleaseVersion(version)
		if err != nil {
			continue
		}
		parsed = append(parsed, parsedVersion)
	}

	return parsed
}
