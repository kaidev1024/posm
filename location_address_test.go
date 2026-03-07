package posm

import (
	"testing"
)

func TestAddressMethods(t *testing.T) {
	var nilAddress *address
	if nilAddress.getCity() != "" {
		t.Fatalf("nil getCity should return empty string")
	}
	if nilAddress.getStreet() != "" {
		t.Fatalf("nil getStreet should return empty string")
	}
	if nilAddress.isCity() {
		t.Fatalf("nil isCity should be false")
	}
	if nilAddress.getAddress() != "" {
		t.Fatalf("nil getAddress should return empty string")
	}

	a := &address{
		HouseNumber: "10",
		Road:        "Market St",
		City:        "San Francisco",
		State:       "CA",
		Postcode:    "94105",
	}
	if got := a.getCity(); got != "San Francisco" {
		t.Fatalf("getCity() = %q", got)
	}
	if got := a.getStreet(); got != "Market St" {
		t.Fatalf("getStreet() = %q", got)
	}
	if a.isCity() {
		t.Fatalf("isCity() should be false when street is present")
	}
	if got := a.getAddress(); got != "10 Market St, San Francisco, CA, 94105" {
		t.Fatalf("getAddress() = %q", got)
	}

	b := &address{Town: "Oakland"}
	if got := b.getCity(); got != "Oakland" {
		t.Fatalf("fallback getCity() = %q", got)
	}
	if !b.isCity() {
		t.Fatalf("isCity() should be true for city-only address")
	}
}

func TestLocationIQResponseMethods(t *testing.T) {
	var nilResp *locationIQResponse
	if nilResp.getPointAddress() != "" || nilResp.getCityAddress() != "" || nilResp.getStreetAddress() != "" {
		t.Fatalf("nil response address getters should return empty strings")
	}
	if nilResp.getStreetSearchText() != "" || nilResp.getCitySearchText() != "" || nilResp.getPlaceID() != "" {
		t.Fatalf("nil response search/id getters should return empty strings")
	}
	if nilResp.isCity() {
		t.Fatalf("nil response isCity should be false")
	}

	lat, lng, err := nilResp.parseCoordinates()
	if err == nil || lat != HEADQUARTER_LAT || lng != HEADQUARTER_LNG {
		t.Fatalf("nil parseCoordinates should return headquarters coordinates and error")
	}

	resp := &locationIQResponse{
		PlaceID:     "123",
		DisplayName: "Display Name",
		Lat:         "37.79",
		Lng:         "-122.39",
		Address: &address{
			HouseNumber: "10",
			Road:        "Market St",
			City:        "San Francisco",
			State:       "CA",
			CountryCode: "us",
		},
	}

	if got := resp.getPointAddress(); got != "10 Market St, San Francisco, CA" {
		t.Fatalf("getPointAddress() = %q", got)
	}
	if got := resp.getCityAddress(); got != "San Francisco, CA" {
		t.Fatalf("getCityAddress() = %q", got)
	}
	if got := resp.getStreetAddress(); got != "Market St, San Francisco, CA" {
		t.Fatalf("getStreetAddress() = %q", got)
	}
	if got := resp.getStreetSearchText(); got != "Market St, San Francisco, CA, us" {
		t.Fatalf("getStreetSearchText() = %q", got)
	}
	if got := resp.getCitySearchText(); got != "San Francisco, CA, us" {
		t.Fatalf("getCitySearchText() = %q", got)
	}

	parsedLat, parsedLng, parseErr := resp.parseCoordinates()
	if parseErr != nil || parsedLat != 37.79 || parsedLng != -122.39 {
		t.Fatalf("parseCoordinates() = (%v, %v, %v)", parsedLat, parsedLng, parseErr)
	}

	if resp.isCity() {
		t.Fatalf("isCity() should be false when street exists")
	}

	cityOnlyResp := &locationIQResponse{Address: &address{City: "San Francisco"}}
	if !cityOnlyResp.isCity() {
		t.Fatalf("isCity() should be true for city-only address")
	}

	// getPlaceID branches
	base := &locationIQResponse{PlaceID: "abc"}
	if got := base.getPlaceID(); got != "abc" {
		t.Fatalf("getPlaceID (empty osm fields) = %q", got)
	}
	withPlace := &locationIQResponse{PlaceID: "123", OsmID: "99", OsmType: "way"}
	if got := withPlace.getPlaceID(); got != "P123" {
		t.Fatalf("getPlaceID (place present) = %q", got)
	}
	fallback := &locationIQResponse{OsmID: "99", OsmType: "way", DisplayName: "Main Road", Lat: "1", Lng: "2"}
	if got := fallback.getPlaceID(); got != "XMain_Road_1_2" {
		t.Fatalf("getPlaceID (fallback) = %q", got)
	}
}
