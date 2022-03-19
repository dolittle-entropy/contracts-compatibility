package main

import (
	"dolittle.io/contracts-compatibility/artifacts"
	"github.com/coreos/go-semver/semver"
)

// ResolveCompatibilityFrom resolves a new Compatibility from a Graph of released artifacts
func ResolveCompatibilityFrom(graph *artifacts.Graph) *Compatibility {
	compatibility := &Compatibility{
		Runtime: make(map[string]map[string]semver.Versions),
		SDKs:    make(map[string]map[string]semver.Versions),
	}

	for _, release := range graph.Runtime {
		sdkCompatibility := make(map[string]semver.Versions)

		for sdk, releases := range graph.SDKs {
			sdkCompatibility[sdk] = resolveCompatibleSDKReleases(release, sdk, releases)
		}

		compatibility.Runtime[release.Version.String()] = sdkCompatibility
	}

	for sdk, releases := range graph.SDKs {
		runtimeCompatibility := make(map[string]semver.Versions)

		for _, release := range releases {
			runtimeCompatibility[release.Version.String()] = resolveCompatibleRuntimeReleases(sdk, release, graph.Runtime)
		}

		compatibility.SDKs[sdk] = runtimeCompatibility
	}

	return compatibility
}

func resolveCompatibleSDKReleases(runtime *artifacts.Release, sdk string, releases artifacts.Releases) semver.Versions {
	compatibleVersions := make(semver.Versions, 0)

	for _, release := range releases {
		if runtime.ContractsVersion.Major != release.ContractsVersion.Major {
			continue
		}
		if runtime.ContractsVersion.Minor < release.ContractsVersion.Minor {
			continue
		}
		if VersionsCompatibilityIsOverridden(runtime.Version, sdk, release.Version) {
			continue
		}

		compatibleVersions = append(compatibleVersions, release.Version)
	}

	semver.Sort(compatibleVersions)
	return compatibleVersions
}

func resolveCompatibleRuntimeReleases(sdk string, release *artifacts.Release, releases artifacts.Releases) semver.Versions {
	compatibleVersions := make(semver.Versions, 0)

	for _, runtime := range releases {
		if release.ContractsVersion.Major != runtime.ContractsVersion.Major {
			continue
		}
		if release.ContractsVersion.Minor > runtime.ContractsVersion.Minor {
			continue
		}
		if VersionsCompatibilityIsOverridden(runtime.Version, sdk, release.Version) {
			continue
		}

		compatibleVersions = append(compatibleVersions, runtime.Version)
	}

	semver.Sort(compatibleVersions)
	return compatibleVersions
}
