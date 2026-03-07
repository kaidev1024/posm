package posm

import (
	"fmt"
	"net/http"
	"strings"
)

const HEADQUARTER_LAT float64 = 37.7955
const HEADQUARTER_LNG float64 = -122.3937

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

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

func GetStreetByText(text string) (*OsmStreet, error) {
	var globalErr error
	location, err := searchText(text)
	if err != nil {
		return nil, fmt.Errorf("searchText error: %w", err)
	}
	lat, lng, err := location.parseCoordinates()
	if err != nil {
		globalErr = fmt.Errorf("parseCoordinates error: %w", err)
	}
	return &OsmStreet{
		PlaceID:     location.getPlaceID(),
		Lat:         lat,
		Lng:         lng,
		DisplayName: location.DisplayName,
		Address:     location.getStreetAddress(),
	}, globalErr
}

func GetCityByText(text string) (*OsmCity, error) {
	var globalErr error
	location, err := searchText(text)
	if err != nil {
		return nil, fmt.Errorf("searchText error: %w", err)
	}
	lat, lng, err := location.parseCoordinates()
	if err != nil {
		globalErr = fmt.Errorf("parseCoordinates error: %w", err)
	}
	return &OsmCity{
		PlaceID:     location.getPlaceID(),
		Lat:         lat,
		Lng:         lng,
		DisplayName: location.DisplayName,
		Address:     location.getCityAddress(),
	}, globalErr
}

func GetPointByOsmTID(tid string) (*OsmPoint, error) {
	var globalErr error
	point, err := lookupByOsmTID(tid)
	if err != nil {
		return nil, fmt.Errorf("lookup error: %w", err)
	}
	lat, lng, err := point.parseCoordinates()
	if err != nil {
		globalErr = fmt.Errorf("parseCoordinates error: %w", err)
	}
	return &OsmPoint{
		PlaceID:          point.getPlaceID(),
		Lat:              lat,
		Lng:              lng,
		DisplayName:      point.DisplayName,
		Address:          point.getPointAddress(),
		StreetSearchText: point.getStreetSearchText(),
		CitySearchText:   point.getCitySearchText(),
	}, globalErr
}

func GetCityByOsmTID(tid string) (*OsmCity, error) {
	var globalErr error
	city, err := lookupByOsmTID(tid)
	if err != nil {
		return nil, fmt.Errorf("lookup error: %w", err)
	}
	lat, lng, err := city.parseCoordinates()
	if err != nil {
		globalErr = fmt.Errorf("parseCoordinates error: %w", err)
	}
	return &OsmCity{
		PlaceID:     city.getPlaceID(),
		Lat:         lat,
		Lng:         lng,
		DisplayName: city.DisplayName,
		Address:     city.getCityAddress(),
	}, globalErr
}

func GetPointsBySearch(text string) ([]*OsmPoint, error) {
	var globalErr error
	locations, err := searchTextMany(text)
	if err != nil {
		if isUnableToGeocode(err) {
			return []*OsmPoint{}, nil
		}
		return nil, fmt.Errorf("searchTextMany error: %w", err)
	}
	points := make([]*OsmPoint, 0)
	normalizedSearch := strings.ToLower(strings.TrimSpace(text))
	if len(normalizedSearch) > 3 {
		normalizedSearch = normalizedSearch[:3]
	}
	seenAddresses := make(map[string]struct{})
	for _, location := range locations {
		point, err := getOsmPointFromLocationIQResponse(&location)
		if err == nil {
			normalizedAddress := strings.ToLower(strings.TrimSpace(point.Address))
			if normalizedSearch != "" && !strings.HasPrefix(normalizedAddress, normalizedSearch) {
				continue
			}
			if _, exists := seenAddresses[normalizedAddress]; exists {
				continue
			}
			seenAddresses[normalizedAddress] = struct{}{}
			points = append(points, point)
		} else {
			// log the error and continue with other results
			globalErr = fmt.Errorf("getOsmPointFromLocationIQResponse error: %w", err)
		}
	}
	return points, globalErr
}

func GetCitiesBySearch(text string) ([]*OsmCity, error) {
	var globalErr error
	locations, err := searchTextMany(text)
	if err != nil {
		return nil, fmt.Errorf("searchTextMany error: %w", err)
	}
	cities := make([]*OsmCity, 0)
	for _, location := range locations {
		if !location.isCity() {
			continue
		}
		city, err := getOsmCityFromLocationIQResponse(&location)
		if err == nil {
			cities = append(cities, city)
		} else {
			// log the error and continue with other results
			globalErr = fmt.Errorf("getOsmCityFromLocationIQResponse error: %w", err)
		}
	}
	return cities, globalErr
}

func GetCitiesByAutocomplete(text string) ([]*OsmCity, error) {
	var globalErr error
	locations, err := autocomplete(text)
	if err != nil {
		if isUnableToGeocode(err) {
			return []*OsmCity{}, nil
		}
		return nil, fmt.Errorf("autocomplete error: %w", err)
	}
	cities := make([]*OsmCity, 0)
	normalizedSearch := strings.ToLower(strings.TrimSpace(text))
	if len(normalizedSearch) > 3 {
		normalizedSearch = normalizedSearch[:3]
	}
	seenAddresses := make(map[string]struct{})
	for _, location := range locations {
		if !location.isCity() {
			continue
		}
		city, err := getOsmCityFromLocationIQResponse(&location)
		if err == nil {
			normalizedAddress := strings.ToLower(strings.TrimSpace(city.Address))
			if normalizedSearch != "" && !strings.HasPrefix(normalizedAddress, normalizedSearch) {
				continue
			}
			if _, exists := seenAddresses[normalizedAddress]; exists {
				continue
			}
			seenAddresses[normalizedAddress] = struct{}{}
			cities = append(cities, city)
		} else {
			// log the error and continue with other results
			globalErr = fmt.Errorf("getOsmCityFromLocationIQResponse error: %w", err)
		}
	}
	return cities, globalErr
}

func getOsmPointFromLocationIQResponse(resp *locationIQResponse) (*OsmPoint, error) {
	var globalErr error
	lat, lng, err := resp.parseCoordinates()
	if err != nil {
		globalErr = fmt.Errorf("parseCoordinates error: %w", err)
	}
	return &OsmPoint{
		PlaceID:          resp.getPlaceID(),
		Lat:              lat,
		Lng:              lng,
		DisplayName:      resp.DisplayName,
		Address:          resp.getPointAddress(),
		StreetSearchText: resp.getStreetSearchText(),
		CitySearchText:   resp.getCitySearchText(),
	}, globalErr
}

func getOsmCityFromLocationIQResponse(resp *locationIQResponse) (*OsmCity, error) {
	var globalErr error
	lat, lng, err := resp.parseCoordinates()
	if err != nil {
		globalErr = fmt.Errorf("parseCoordinates error: %w", err)
	}

	return &OsmCity{
		PlaceID:     resp.getPlaceID(),
		Lat:         lat,
		Lng:         lng,
		DisplayName: resp.DisplayName,
		Address:     resp.getCityAddress(),
	}, globalErr
}

func isUnableToGeocode(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unable to geocode")
}
