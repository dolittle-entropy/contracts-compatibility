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
			sdkCompatibility[sdk] = resolveCompatibleReleases(release.ContractsVersion, releases, true)
		}

		compatibility.Runtime[release.Version.String()] = sdkCompatibility
	}

	for sdk, releases := range graph.SDKs {
		runtimeCompatibility := make(map[string]semver.Versions)

		for _, release := range releases {
			runtimeCompatibility[release.Version.String()] = resolveCompatibleReleases(release.ContractsVersion, graph.Runtime, false)
		}

		compatibility.SDKs[sdk] = runtimeCompatibility
	}

	return compatibility
}

func resolveCompatibleReleases(contracts *semver.Version, releases artifacts.Releases, contractsShouldBeGreater bool) semver.Versions {
	compatibleVersions := make(semver.Versions, 0)

	for _, release := range releases {
		if contracts.Major != release.ContractsVersion.Major {
			continue
		}
		if contractsShouldBeGreater && contracts.Minor < release.ContractsVersion.Minor {
			continue
		}
		if !contractsShouldBeGreater && contracts.Minor > release.ContractsVersion.Minor {
			continue
		}

		compatibleVersions = append(compatibleVersions, release.Version)
	}

	semver.Sort(compatibleVersions)
	return compatibleVersions
}
