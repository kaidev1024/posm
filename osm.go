package posm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var searchClient *Client
var autoCompleteClient *Client
var lookupClient *Client
var locationIQAccessToken string

// searchText search for OSM location by text, returns the first result
func searchText(query string) (*locationIQResponse, error) {
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

// searchTextMany search for OSM location by text, return all results
func searchTextMany(query string) ([]locationIQResponse, error) {
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

// autocomplete search for OSM location by text, return all results
func autocomplete(query string) ([]locationIQResponse, error) {
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

// lookup search for OSM location by OSM IDs
func lookup(osmTID string) (*locationIQResponse, error) {
	params := url.Values{}
	params.Set("key", locationIQAccessToken)
	params.Set("format", "json")
	params.Set("osm_ids", osmTID)
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
