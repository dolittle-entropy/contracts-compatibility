package main

import (
	"github.com/coreos/go-semver/semver"
)

// Override represents a structure that describes a compatibility breaking checkpoint.
// This means that any versions of a (Runtime, SDK) pair that cross this boundary is considered incompatible
type Override struct {
	RuntimeVersion *semver.Version
	SDKVersions    map[string]*semver.Version
}

// Overrides defines the known breaking release pairs that is not captured in the Contracts versions
var Overrides = []Override{
	{
		// The breaking change in ReverseCall ping-pong behaviour introduced in v6
		RuntimeVersion: semver.New("6.0.0"),
		SDKVersions: map[string]*semver.Version{
			"DotNET":     semver.New("9.0.0"),
			"JavaScript": semver.New("15.0.0"),
		},
	},
}

// VersionsCompatibilityIsOverridden checks whether a pair of (Runtime, SDK) versions have been marked explicitly as not compatible through the Overrides
func VersionsCompatibilityIsOverridden(runtimeVersion *semver.Version, sdk string, sdkVersion *semver.Version) bool {
	for _, override := range Overrides {
		overrideSDKVersion, isOverridden := override.SDKVersions[sdk]
		if !isOverridden {
			continue
		}

		if !runtimeVersion.LessThan(*override.RuntimeVersion) && sdkVersion.LessThan(*overrideSDKVersion) {
			return true
		}

		if runtimeVersion.LessThan(*override.RuntimeVersion) && !sdkVersion.LessThan(*overrideSDKVersion) {
			return true

		}
	}

	return false
}
