package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DoJSON sends and HTTP request and decodes the JSON-encoded response in the value pointed to by v.
func DoJSON(req *http.Request, v any) error {
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get JSON data: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %v", response.Status)
	}

	err = json.NewDecoder(response.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("failed to decode received JSON data: %w", err)
	}

	return nil
}

// GetJSON issues a GET to the specified URl and decodes the JSON-encoded response in the value pointed to by v.
func GetJSON(url string, v any) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Accept", "application/json")

	return DoJSON(request, v)
}
