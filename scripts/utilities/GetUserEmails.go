package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetUserEmails(accessToken string) ([]map[string]interface{}, error) {
	var bodyBytes []byte

	getUserParams := map[string]interface{}{
		"URL":    "https://api.github.com/user/emails",
		"Method": "GET",
		"Headers": map[string]string{
			"Accept":        "application/json",
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("token %v", accessToken),
		},
	}

	req, err := http.NewRequest(getUserParams["Method"].(string), getUserParams["URL"].(string), nil)

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	for key, value := range getUserParams["Headers"].(map[string]string) {
		req.Header.Set(key, value)
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	defer resp.Body.Close()

	bodyBytes, err = io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var emails []map[string]interface{}

	err = json.Unmarshal(bodyBytes, &emails)

	if err != nil {
		return nil, fmt.Errorf("failed to get user emails: %v", err)
	}

	return emails, nil
}
