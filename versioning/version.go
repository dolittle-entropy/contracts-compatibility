package versioning

import (
	"errors"
	"github.com/coreos/go-semver/semver"
)

// NewReleaseVersion parses the version string and returns the value if it is not a prerelease
func NewReleaseVersion(version string) (*semver.Version, error) {
	parsed, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	if parsed.PreRelease != "" {
		return nil, errors.New("version is a prerelease")
	}

	return parsed, nil
}
