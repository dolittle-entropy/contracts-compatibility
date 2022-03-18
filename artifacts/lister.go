package artifacts

import (
	"fmt"
	"github.com/coreos/go-semver/semver"
	"sync"
)

// Lister defines a system that can list released versions of an artifact
type Lister interface {
	ListReleasedVersions() (semver.Versions, error)
}

// Resolver defines a system that can resolve the Contracts dependency for an artifact for a specific version
type Resolver interface {
	ResolveDependencyForVersion(version *semver.Version) (*semver.Version, error)
}

// ReleaseListResolver represents a system that can list all releases and resolves their Contracts dependency
type ReleaseListResolver struct {
	lister   Lister
	resolver Resolver
}

// NewReleaseListResolver creates a new ReleaseListResolver using the given Lister and Resolver
func NewReleaseListResolver(lister Lister, resolver Resolver) *ReleaseListResolver {
	return &ReleaseListResolver{
		lister:   lister,
		resolver: resolver,
	}
}

// ListAndResolve lists all releases and resolves their Contracts dependency
func (listResolver *ReleaseListResolver) ListAndResolve() (Releases, error) {
	releasedVersions, err := listResolver.lister.ListReleasedVersions()
	if err != nil {
		return nil, fmt.Errorf("could not list released versions: %w", err)
	}
	semver.Sort(releasedVersions)

	ch := make(chan *Release, len(releasedVersions))
	wg := sync.WaitGroup{}

	for _, releasedVersion := range releasedVersions {
		wg.Add(1)
		go func(releasedVersion *semver.Version) {
			defer wg.Done()
			fmt.Println("Resolving dependency for version", releasedVersion.String())
			contractsVersion, err := listResolver.resolver.ResolveDependencyForVersion(releasedVersion)
			if err != nil {
				return
			}
			ch <- &Release{
				Version:          releasedVersion,
				ContractsVersion: contractsVersion,
			}
		}(releasedVersion)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	releases := make(Releases, 0)
	for release := range ch {
		releases = append(releases, release)
	}
	return releases, nil
}
