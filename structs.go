package posm

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type OsmType uint8

const (
	OsmTypeNone OsmType = iota
	OsmTypeNode
	OsmTypeWay
	OsmTypeRelation
)

const INVALID_LAT float64 = 200.0
const INVALID_LNG float64 = 200.0
const INVALID_OSM_ID int64 = 0
const HEADQUARTER_LAT float64 = 37.7955
const HEADQUARTER_LNG float64 = -122.3937

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// Address contains address fields specific to OpenStreetMap
type Address struct {
	HouseNumber   string `json:"house_number,omitempty"`
	Road          string `json:"road,omitempty"`
	Pedestrian    string `json:"pedestrian,omitempty"`
	Footway       string `json:"footway,omitempty"`
	Cycleway      string `json:"cycleway,omitempty"`
	Highway       string `json:"highway,omitempty"`
	Path          string `json:"path,omitempty"`
	Suburb        string `json:"suburb,omitempty"`
	City          string `json:"city,omitempty"`
	Town          string `json:"town,omitempty"`
	Village       string `json:"village,omitempty"`
	Hamlet        string `json:"hamlet,omitempty"`
	County        string `json:"county,omitempty"`
	Country       string `json:"country,omitempty"`
	CountryCode   string `json:"country_code,omitempty"`
	State         string `json:"state,omitempty"`
	StateDistrict string `json:"state_district,omitempty"`
	Postcode      string `json:"postcode,omitempty"`
}

// getCity checks different fields for the city name
func (a *Address) getCity() string {
	if a == nil {
		return ""
	}
	var city string
	if a.City != "" {
		city = a.City
	} else if a.Town != "" {
		city = a.Town
	} else if a.Village != "" {
		city = a.Village
	} else if a.Hamlet != "" {
		city = a.Hamlet
	}
	return city
}

// getStreet checks different fields for the street name
func (a *Address) getStreet() string {
	if a == nil {
		return ""
	}
	var street string
	if a.Road != "" {
		street = a.Road
	} else if a.Pedestrian != "" {
		street = a.Pedestrian
	} else if a.Path != "" {
		street = a.Path
	} else if a.Cycleway != "" {
		street = a.Cycleway
	} else if a.Footway != "" {
		street = a.Footway
	} else if a.Highway != "" {
		street = a.Highway
	}
	return street
}

func (a *Address) getAddress() string {
	if a == nil {
		return ""
	}
	address := a.getCity()
	street := a.getStreet()
	if street != "" {
		address = fmt.Sprintf("%s, %s", street, address)
		if a.HouseNumber != "" {
			address = fmt.Sprintf("%s %s", a.HouseNumber, address)
		}
	}
	if a.State != "" {
		address = fmt.Sprintf("%s, %s", address, a.State)
	}
	if a.Postcode != "" {
		address = fmt.Sprintf("%s, %s", address, a.Postcode)
	}
	return address
}

func truncateAtFirstComma(s string) string {
	if idx := strings.Index(s, ","); idx != -1 {
		return s[:idx]
	}
	return s // no comma found
}

type LocationIQResponse struct {
	OsmID       string   `json:"osm_id"`
	OsmType     string   `json:"osm_type"`
	DisplayName string   `json:"display_name"`
	Lat         string   `json:"lat"`
	Lng         string   `json:"lon"`
	Address     *Address `json:"address"`
}

func (lr *LocationIQResponse) GetPointAddress() string {
	if lr == nil {
		return ""
	}
	name := truncateAtFirstComma(lr.DisplayName)
	address := lr.Address.getAddress()
	if idx := strings.Index(address, name); idx == -1 {
		address = fmt.Sprintf("%s, %s", name, address)
	}
	return address
}

func (lr *LocationIQResponse) GetCityAddress() string {
	if lr == nil {
		return ""
	}
	address := lr.Address
	city := address.getCity()
	if city == "" {
		city = address.County
	}
	return fmt.Sprintf("%s, %s", city, address.State)
}

func (lr *LocationIQResponse) GetStreetAddress() string {
	if lr == nil {
		return ""
	}
	address := lr.Address
	return fmt.Sprintf("%s, %s, %s", address.getStreet(), address.getCity(), address.State)
}

func (lr *LocationIQResponse) GetStreetAddressForSearch() string {
	if lr == nil {
		return ""
	}
	address := lr.Address
	street := address.getStreet()
	if street == "" {
		return ""
	}
	return fmt.Sprintf("%s, %s, %s, %s", street, address.getCity(), address.State, address.CountryCode)
}

func (lr *LocationIQResponse) GetCityAddressForSearch() string {
	if lr == nil {
		return ""
	}
	address := lr.Address
	city := address.getCity()
	if city == "" {
		return ""
	}
	return fmt.Sprintf("%s, %s, %s", city, address.State, address.CountryCode)
}

func (lr *LocationIQResponse) GetLat() (float64, error) {
	if lr == nil {
		return INVALID_LAT, fmt.Errorf("empty location")
	}
	lat, err := strconv.ParseFloat(lr.Lat, 64)
	if err != nil {
		return INVALID_LAT, err
	}
	return lat, nil
}

func (lr *LocationIQResponse) GetLng() (float64, error) {
	if lr == nil {
		return INVALID_LNG, fmt.Errorf("empty location")
	}
	lng, err := strconv.ParseFloat(lr.Lng, 64)
	if err != nil {
		return INVALID_LNG, err
	}
	return lng, nil
}

func (lr *LocationIQResponse) GetOsmID() (int64, error) {
	if lr == nil {
		return 0, fmt.Errorf("empty location")
	}
	osmID, err := strconv.ParseInt(lr.OsmID, 10, 64)
	if err != nil {
		return INVALID_OSM_ID, err
	}
	return osmID, nil
}

func (lr *LocationIQResponse) GetOsmType() OsmType {
	if lr == nil {
		return OsmTypeNone
	}
	if lr.OsmType == "node" {
		return OsmTypeNode
	}
	if lr.OsmType == "way" {
		return OsmTypeWay
	}
	if lr.OsmType == "relation" {
		return OsmTypeRelation
	}
	return OsmTypeNone
}
