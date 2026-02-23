package posm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

var searchClient *Client
var autoCompleteClient *Client
var lookupClient *Client
var locationIQAccessToken string

// NewClient creates a new LocationIQ client
func Init(accessToken string) {
	locationIQAccessToken = accessToken
	searchClient = &Client{
		BaseURL:    "https://us1.locationiq.com/v1/search",
		HTTPClient: &http.Client{},
	}
	autoCompleteClient = &Client{
		BaseURL:    "https://api.locationiq.com/v1/autocomplete",
		HTTPClient: &http.Client{},
	}
	lookupClient = &Client{
		BaseURL:    "https://us1.locationiq.com/v1/lookup",
		HTTPClient: &http.Client{},
	}
}

// SearchText search for OSM location by text, returns the first result
func SearchText(query string) (*locationIQResponse, error) {
	params := url.Values{}
	params.Set("key", locationIQAccessToken)
	params.Set("format", "json")
	params.Set("addressdetails", "1")
	params.Set("q", query)
	reqURL := searchClient.BaseURL + "?" + params.Encode()
	resp, err := searchClient.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}
	var results []locationIQResponse
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

// SearchTextMany search for OSM location by text, return all results
func SearchTextMany(query string) ([]locationIQResponse, error) {
	params := url.Values{}
	params.Set("key", locationIQAccessToken)
	params.Set("format", "json")
	params.Set("addressdetails", "1")
	params.Set("q", query)
	reqURL := searchClient.BaseURL + "?" + params.Encode()
	resp, err := searchClient.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}
	var results []locationIQResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return results, nil
}

// Autocomplete search for OSM location by text, return all results
func Autocomplete(query string) ([]locationIQResponse, error) {
	params := url.Values{}
	params.Set("key", locationIQAccessToken)
	params.Set("format", "json")
	params.Set("dedupe", "1")
	params.Set("limit", "10")
	params.Set("q", query)
	reqURL := autoCompleteClient.BaseURL + "?" + params.Encode()
	resp, err := autoCompleteClient.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}
	var results []locationIQResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return results, nil
}

// Lookup search for OSM location by OSM IDs
func Lookup(osmID string) (*locationIQResponse, error) {
	params := url.Values{}
	params.Set("key", locationIQAccessToken)
	params.Set("format", "json")
	params.Set("osm_ids", osmID)
	reqURL := lookupClient.BaseURL + "?" + params.Encode()
	resp, err := lookupClient.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}
	var results []locationIQResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no results found")
	}
	return &results[0], nil
}
