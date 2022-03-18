package docker

import (
	"archive/tar"
	"compress/gzip"
	"dolittle.io/contracts-compatibility/http"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"io"
	goHTTP "net/http"
)

// FileDependencyResolver defines a system that can resolve the version of a dependency from a file
type FileDependencyResolver interface {
	ResolveDependencyFromContents(contents []byte) (*semver.Version, error)
}

// DependencyResolver represents a system that can resolve the version of a dependency from a file in a Docker image
type DependencyResolver struct {
	authToken AuthToken
	imageName string
	fileNames []string
	resolver  FileDependencyResolver
}

// NewDependencyResolverFor creates a new DependencyResolver for a Docker image with the given name, files to parse and resolver using the provided authentication token
func NewDependencyResolverFor(token AuthToken, image string, resolver FileDependencyResolver, files ...string) *DependencyResolver {
	return &DependencyResolver{
		authToken: token,
		imageName: image,
		fileNames: files,
		resolver:  resolver,
	}
}

// ResolveDependencyForVersion resolves the dependency version for the given package version
func (resolver *DependencyResolver) ResolveDependencyForVersion(version *semver.Version) (*semver.Version, error) {
	layers, err := resolver.getDockerImageLayers(version)
	if err != nil {
		return nil, err
	}

	for _, layer := range resolver.sortMostLikelyLayers(layers) {
		fileContents, err := resolver.getFileInLayer(layer)
		if err == io.EOF {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("could not find file in layer: %w", err)
		}

		return resolver.resolver.ResolveDependencyFromContents(fileContents)
	}

	return nil, fmt.Errorf("no dependency found for %v", resolver.imageName)
}

func (resolver *DependencyResolver) getDockerImageLayers(version *semver.Version) ([]dockerLayer, error) {
	request, err := CreateAuthenticatedGETRequestTo(resolver.authToken, "https://registry-1.docker.io/v2/"+resolver.imageName+"/manifests/"+version.String())
	if err != nil {
		return nil, fmt.Errorf("could not create Docker manifest request: %w", err)
	}
	request.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	request.Header.Add("Accept", "application/vnd.docker.distribution.manifest.list.v2+json")

	info := dockerManifest{}
	if err := http.DoJSON(request, &info); err != nil {
		return nil, fmt.Errorf("could not get Docker manifest: %w", err)
	}

	return info.Layers, nil
}

func (resolver *DependencyResolver) sortMostLikelyLayers(layers []dockerLayer) []dockerLayer {
	sorted := make([]dockerLayer, 0)
	for i := len(layers) - 1; i >= 0; i-- {
		layer := layers[i]
		if layer.Size < minLikelyLayerSize {
			continue
		}
		sorted = append(sorted, layer)
	}
	return sorted
}

func (resolver *DependencyResolver) getFileInLayer(layer dockerLayer) ([]byte, error) {
	request, err := CreateAuthenticatedGETRequestTo(resolver.authToken, "https://registry-1.docker.io/v2/"+resolver.imageName+"/blobs/"+layer.Digest)
	if err != nil {
		return nil, fmt.Errorf("could not create Docker blob request: %w", err)
	}
	request.Header.Set("Accept", "application/octet-stream")

	response, err := goHTTP.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not GET blob data: %w", err)
	}
	defer response.Body.Close()

	decompressReader, err := gzip.NewReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("could not decompress blob data: %w", err)
	}
	defer decompressReader.Close()

	archiveReader := tar.NewReader(decompressReader)
	for {
		fileHeader, err := archiveReader.Next()
		if err != nil {
			return nil, err
		}

		for _, fileName := range resolver.fileNames {
			if !fileHeader.FileInfo().IsDir() && fileHeader.Name == fileName {
				return io.ReadAll(archiveReader)
			}
		}
	}
}

const minLikelyLayerSize = 2048
