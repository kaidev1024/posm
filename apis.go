package posm

import (
	"fmt"
	"net/http"
)

const (
	OsmTypeNone OsmType = iota
	OsmTypeNode
	OsmTypeWay
	OsmTypeRelation
)

const INVALID_LAT float64 = 999.0
const INVALID_LNG float64 = 999.0
const INVALID_OSM_ID int64 = 0
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
	osmID, err := location.getOsmID()
	if err != nil {
		globalErr = fmt.Errorf("getOsmID error: %w", err)
	}
	osmType := location.getOsmType()
	lat, err := location.getLat()
	if err != nil {
		globalErr = fmt.Errorf("getLat error: %w", err)
	}
	lng, err := location.getLng()
	if err != nil {
		globalErr = fmt.Errorf("getLng error: %w", err)
	}
	return &OsmStreet{
		OsmID:       osmID,
		OsmType:     osmType,
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
	osmID, err := location.getOsmID()
	if err != nil {
		globalErr = fmt.Errorf("getOsmID error: %w", err)
	}
	osmType := location.getOsmType()
	lat, err := location.getLat()
	if err != nil {
		globalErr = fmt.Errorf("getLat error: %w", err)
	}
	lng, err := location.getLng()
	if err != nil {
		globalErr = fmt.Errorf("getLng error: %w", err)
	}
	return &OsmCity{
		OsmID:       osmID,
		OsmType:     osmType,
		Lat:         lat,
		Lng:         lng,
		DisplayName: location.DisplayName,
		Address:     location.getCityAddress(),
	}, globalErr
}

func GetPointByTID(tid string) (*OsmPoint, error) {
	var globalErr error
	location, err := lookup(tid)
	if err != nil {
		return nil, fmt.Errorf("lookup error: %w", err)
	}
	osmID, err := location.getOsmID()
	if err != nil {
		globalErr = fmt.Errorf("getOsmID error: %w", err)
	}
	osmType := location.getOsmType()
	lat, err := location.getLat()
	if err != nil {
		globalErr = fmt.Errorf("getLat error: %w", err)
	}
	lng, err := location.getLng()
	if err != nil {
		globalErr = fmt.Errorf("getLng error: %w", err)
	}
	return &OsmPoint{
		OsmID:            osmID,
		OsmType:          osmType,
		Lat:              lat,
		Lng:              lng,
		DisplayName:      location.DisplayName,
		Address:          location.getPointAddress(),
		StreetSearchText: location.getStreetSearchText(),
		CitySearchText:   location.getCitySearchText(),
	}, globalErr
}

func GetPointsBySearch(text string) ([]*OsmPoint, error) {
	var globalErr error
	locations, err := searchTextMany(text)
	if err != nil {
		return nil, fmt.Errorf("searchTextMany error: %w", err)
	}
	points := make([]*OsmPoint, 0)
	for _, location := range locations {
		point, err := getOsmPointFromLocationIQResponse(&location)
		if err == nil {
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
		return nil, fmt.Errorf("autocomplete error: %w", err)
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

func getOsmPointFromLocationIQResponse(resp *locationIQResponse) (*OsmPoint, error) {
	var globalErr error
	osmID, err := resp.getOsmID()
	if err != nil {
		globalErr = fmt.Errorf("getOsmID error: %w", err)
	}
	osmType := resp.getOsmType()
	lat, err := resp.getLat()
	if err != nil {
		globalErr = fmt.Errorf("getLat error: %w", err)
	}
	lng, err := resp.getLng()
	if err != nil {
		globalErr = fmt.Errorf("getLng error: %w", err)
	}
	return &OsmPoint{
		OsmID:            osmID,
		OsmType:          osmType,
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
	osmID, err := resp.getOsmID()
	if err != nil {
		globalErr = fmt.Errorf("getOsmID error: %w", err)
	}
	osmType := resp.getOsmType()
	lat, err := resp.getLat()
	if err != nil {
		globalErr = fmt.Errorf("getLat error: %w", err)
	}
	lng, err := resp.getLng()
	if err != nil {
		globalErr = fmt.Errorf("getLng error: %w", err)
	}
	return &OsmCity{
		OsmID:       osmID,
		OsmType:     osmType,
		Lat:         lat,
		Lng:         lng,
		DisplayName: resp.DisplayName,
		Address:     resp.getCityAddress(),
	}, globalErr
}
