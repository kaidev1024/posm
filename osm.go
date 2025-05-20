package posm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var client *Client
var locationIQAccessToken string

// NewClient creates a new LocationIQ client
func Init(accessToken string) {
	locationIQAccessToken = accessToken
	client = &Client{
		BaseURL:    "https://us1.locationiq.com/v1/search.php",
		HTTPClient: &http.Client{},
	}
}

// SearchText search for OSM location by text
func SearchText(query string) (*LocationIQResponse, error) {
	params := url.Values{}
	params.Set("key", locationIQAccessToken)
	params.Set("format", "json")
	params.Set("addressdetails", "1")
	params.Set("q", query)
	reqURL := client.BaseURL + "?" + params.Encode()
	resp, err := client.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}
	var results []LocationIQResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	for _, result := range results {
		if result.Address.getCity() != "" {
			return &result, nil
		}
	}
	return &results[0], nil
}
