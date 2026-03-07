package posm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupClientsForServer(server *httptest.Server) {
	locationIQAccessToken = "test-key"
	searchClient = &Client{BaseURL: server.URL + "/search", HTTPClient: server.Client()}
	autoCompleteClient = &Client{BaseURL: server.URL + "/autocomplete", HTTPClient: server.Client()}
	lookupClient = &Client{BaseURL: server.URL + "/lookup", HTTPClient: server.Client()}
}

func TestInit(t *testing.T) {
	Init("abc123")
	if locationIQAccessToken != "abc123" {
		t.Fatalf("Init did not set access token")
	}
	if searchClient == nil || autoCompleteClient == nil || lookupClient == nil {
		t.Fatalf("Init did not initialize all clients")
	}
}

func TestOsmHTTPFunctions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/search":
			q := r.URL.Query().Get("q")
			switch q {
			case "pick-city":
				_, _ = fmt.Fprint(w, `[
					{"place_id":"1","display_name":"No city","lat":"1","lon":"2","address":{"road":"Any Road"}},
					{"place_id":"2","display_name":"With city","lat":"3","lon":"4","address":{"city":"SF"}}
				]`)
			case "many":
				_, _ = fmt.Fprint(w, `[ {"place_id":"11","display_name":"Result","lat":"1","lon":"2","address":{"city":"SF"}} ]`)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		case "/autocomplete":
			if r.URL.Query().Get("q") == "boom" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = fmt.Fprint(w, `[ {"place_id":"21","display_name":"City","lat":"1","lon":"2","address":{"city":"San Jose","state":"CA","country_code":"us"}} ]`)
		case "/lookup":
			if r.URL.Query().Get("osm_ids") == "EMPTY" {
				_, _ = fmt.Fprint(w, `[]`)
				return
			}
			_, _ = fmt.Fprint(w, `[ {"place_id":"31","display_name":"Point","lat":"10","lon":"20","address":{"road":"Main","city":"SF","state":"CA","country_code":"us"}} ]`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	setupClientsForServer(server)

	// searchText
	resp, err := searchText("pick-city")
	if err != nil || resp.PlaceID != "2" {
		t.Fatalf("searchText failed: resp=%+v err=%v", resp, err)
	}

	// searchTextMany (404 => empty)
	results, err := searchTextMany("unknown")
	if err != nil || len(results) != 0 {
		t.Fatalf("searchTextMany 404 handling failed: len=%d err=%v", len(results), err)
	}

	// autocomplete non-200
	_, err = autocomplete("boom")
	if err == nil {
		t.Fatalf("autocomplete should fail on non-200")
	}

	// lookupByOsmTID empty result
	_, err = lookupByOsmTID("EMPTY")
	if err == nil {
		t.Fatalf("lookupByOsmTID should fail on empty results")
	}
}

func TestAPIFunctions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/search":
			q := r.URL.Query().Get("q")
			switch q {
			case "street":
				_, _ = fmt.Fprint(w, `[
					{"place_id":"s1","display_name":"Street","lat":"37.1","lon":"-122.1","address":{"road":"Market St","city":"San Francisco","state":"CA"}}
				]`)
			case "city":
				_, _ = fmt.Fprint(w, `[
					{"place_id":"c1","display_name":"City","lat":"bad","lon":"-122.1","address":{"city":"San Francisco","state":"CA"}}
				]`)
			case "points":
				_, _ = fmt.Fprint(w, `[
					{"place_id":"p1","display_name":"One","lat":"10","lon":"20","address":{"road":"Main St","city":"San Jose","state":"CA","country_code":"us"}},
					{"place_id":"p2","display_name":"Dup","lat":"11","lon":"21","address":{"road":"Main St","city":"San Jose","state":"CA","country_code":"us"}},
					{"place_id":"p3","display_name":"Bad","lat":"oops","lon":"21","address":{"road":"Main St","city":"San Jose","state":"CA","country_code":"us"}}
				]`)
			case "cities":
				_, _ = fmt.Fprint(w, `[
					{"place_id":"k1","display_name":"City1","lat":"1","lon":"2","address":{"city":"San Jose","state":"CA","country_code":"us"}},
					{"place_id":"k2","display_name":"NotCity","lat":"1","lon":"2","address":{"road":"Road","city":"San Jose","state":"CA","country_code":"us"}}
				]`)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		case "/autocomplete":
			if r.URL.Query().Get("q") == "none" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			_, _ = fmt.Fprint(w, `[
				{"place_id":"a1","display_name":"A","lat":"1","lon":"2","address":{"city":"San Mateo","state":"CA","country_code":"us"}},
				{"place_id":"a2","display_name":"A2","lat":"1","lon":"2","address":{"city":"San Mateo","state":"CA","country_code":"us"}}
			]`)
		case "/lookup":
			_, _ = fmt.Fprint(w, `[
				{"place_id":"l1","display_name":"Lookup","lat":"37.7","lon":"-122.4","address":{"road":"Howard St","city":"San Francisco","state":"CA","country_code":"us"}}
			]`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	setupClientsForServer(server)

	street, err := GetStreetBySearch("street")
	if err != nil || street.PlaceID == "" || street.Address == "" {
		t.Fatalf("GetStreetBySearch failed: street=%+v err=%v", street, err)
	}

	city, err := GetCityBySearch("city")
	if err == nil || city == nil {
		t.Fatalf("GetCityBySearch should return city with parse error, got city=%+v err=%v", city, err)
	}

	point, err := GetPointByLookup("W1")
	if err != nil || point.PlaceID == "" || point.StreetSearchText == "" || point.CitySearchText == "" {
		t.Fatalf("GetPointByLookup failed: point=%+v err=%v", point, err)
	}

	cityByID, err := GetCityByLookup("W1")
	if err != nil || cityByID.PlaceID == "" || cityByID.Address == "" {
		t.Fatalf("GetCityByLookup failed: city=%+v err=%v", cityByID, err)
	}

	points, err := GetPointsBySearch("san")
	if len(points) != 1 {
		t.Fatalf("GetPointsBySearch should dedupe/filter, got %d", len(points))
	}
	if err == nil {
		t.Fatalf("GetPointsBySearch should return aggregated error when one item is invalid")
	}

	cities, err := GetCitiesBySearch("cities")
	if err != nil || len(cities) != 1 {
		t.Fatalf("GetCitiesBySearch failed: len=%d err=%v", len(cities), err)
	}

	autoCities, err := GetCitiesByAutocomplete("san")
	if err != nil || len(autoCities) != 1 {
		t.Fatalf("GetCitiesByAutocomplete failed: len=%d err=%v", len(autoCities), err)
	}

	emptyAutoCities, err := GetCitiesByAutocomplete("none")
	if err != nil || len(emptyAutoCities) != 0 {
		t.Fatalf("GetCitiesByAutocomplete should return empty slice on 404: len=%d err=%v", len(emptyAutoCities), err)
	}
}

func TestConvertersAndErrorHelpers(t *testing.T) {
	point, err := getOsmPointFromLocationIQResponse(&locationIQResponse{
		PlaceID:     "1",
		DisplayName: "D",
		Lat:         "bad",
		Lng:         "2",
		Address:     &address{Road: "Road", City: "City", State: "ST", CountryCode: "us"},
	})
	if err == nil || point == nil {
		t.Fatalf("getOsmPointFromLocationIQResponse should return point and parse error")
	}

	city, err := getOsmCityFromLocationIQResponse(&locationIQResponse{
		PlaceID:     "2",
		DisplayName: "D2",
		Lat:         "1",
		Lng:         "2",
		Address:     &address{City: "City", State: "ST", CountryCode: "us"},
	})
	if err != nil || city == nil || city.Address == "" {
		t.Fatalf("getOsmCityFromLocationIQResponse failed: city=%+v err=%v", city, err)
	}

	if !isUnableToGeocode(fmt.Errorf("Unable to Geocode this input")) {
		t.Fatalf("isUnableToGeocode should detect error message")
	}
	if isUnableToGeocode(fmt.Errorf("different error")) {
		t.Fatalf("isUnableToGeocode should return false for unrelated errors")
	}
}
