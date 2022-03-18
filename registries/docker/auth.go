package docker

import (
	"dolittle.io/contracts-compatibility/http"
	"fmt"
	goHTTP "net/http"
)

// AuthToken represents an authentication token to use with the Docker Hub APIs
type AuthToken = string

// GetAuthTokenFor gets an authentication token to use with the Docker Hub APIs for the given image repository
func GetAuthTokenFor(image string) (AuthToken, error) {
	token := dockerAuthToken{}
	if err := http.GetJSON("https://auth.docker.io/token?service=registry.docker.io&scope=repository:"+image+":pull", &token); err != nil {
		return "", fmt.Errorf("failed to get Docker Hub authentication token: %w", err)
	}

	return token.Token, nil
}

// CreateAuthenticatedGETRequestTo creates an HTTP Request to the given URL that is authenticated using the provided token
func CreateAuthenticatedGETRequestTo(token AuthToken, url string) (*goHTTP.Request, error) {
	request, err := goHTTP.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/json")
	return request, nil
}
