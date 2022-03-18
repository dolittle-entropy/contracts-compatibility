package artifacts

import "github.com/coreos/go-semver/semver"

// Release defines a released version of an artifact with a dependency on Contracts
type Release struct {
	Version          *semver.Version `json:"version"`
	ContractsVersion *semver.Version `json:"contracts"`
}

// Releases defines a slice of Release
type Releases = []*Release
