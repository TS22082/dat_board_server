package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPRequestParams holds parameters for the HTTP request
type HTTPRequestParams struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    interface{}
}

// MakeHTTPRequest makes an HTTP request with the given parameters and returns the response body, status code, and error
func MakeHTTPRequest(params HTTPRequestParams) (map[string]interface{}, int, error) {
	// Marshal the body to JSON if it's not nil
	var bodyBytes []byte
	var err error
	if params.Body != nil {
		bodyBytes, err = json.Marshal(params.Body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal body: %v", err)
		}
	}

	// Create a new HTTP request
	req, err := http.NewRequest(params.Method, params.URL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response body: %v", err)
	}

	// Unmarshal the response body into a map
	var result map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return result, resp.StatusCode, nil
}
