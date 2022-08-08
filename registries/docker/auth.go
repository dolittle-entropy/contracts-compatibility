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
	request, err := createAuthTokenRequest(image)
	if err != nil {
		return "", err
	}
	return requestAuthTokenUsing(request)
}

// GetAuthenticatedUserAuthTokenFor gets an authentication token to use with the Docker Hub APIs for the given image repository using the supplied username and password
func GetAuthenticatedUserAuthTokenFor(image, username, password string) (AuthToken, error) {
	request, err := createAuthTokenRequest(image)
	if err != nil {
		return "", err
	}
	request.SetBasicAuth(username, password)
	return requestAuthTokenUsing(request)
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

func createAuthTokenRequest(image string) (*goHTTP.Request, error) {
	return goHTTP.NewRequest("GET", "https://auth.docker.io/token?service=registry.docker.io&scope=repository:"+image+":pull", nil)
}

func requestAuthTokenUsing(request *goHTTP.Request) (AuthToken, error) {
	token := dockerAuthToken{}
	if err := http.DoJSON(request, &token); err != nil {
		return "", fmt.Errorf("failed to get Docker Hub authentication token: %w", err)
	}
	return token.Token, nil
}
