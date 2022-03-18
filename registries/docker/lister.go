package docker

import (
	"dolittle.io/contracts-compatibility/http"
	"dolittle.io/contracts-compatibility/versioning"
	"fmt"
	"github.com/coreos/go-semver/semver"
)

// ReleaseLister represents a system that can list released versions of a Docker image based on tags
type ReleaseLister struct {
	authToken AuthToken
	imageName string
}

// NewReleaseListerFor creates a new ReleaseLister for a Docker image with the given name using the provided authentication token
func NewReleaseListerFor(token AuthToken, image string) *ReleaseLister {
	return &ReleaseLister{
		authToken: token,
		imageName: image,
	}
}

// ListReleasedVersions lists the released version of the Docker image from its tags
func (lister *ReleaseLister) ListReleasedVersions() (semver.Versions, error) {
	request, err := CreateAuthenticatedGETRequestTo(lister.authToken, "https://registry-1.docker.io/v2/"+lister.imageName+"/tags/list")
	if err != nil {
		return nil, fmt.Errorf("could not create Docker tag list request: %w", err)
	}

	info := dockerTagList{}
	if err := http.DoJSON(request, &info); err != nil {
		return nil, fmt.Errorf("could not get Docker tag list: %w", err)
	}

	return versioning.NewReleaseVersions(info.Tags), nil
}
