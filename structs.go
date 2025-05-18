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

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// Address contains address fields specific to OpenStreetMap
type Address struct {
	HouseNumber   string `json:"house_number"`
	Road          string `json:"road"`
	Pedestrian    string `json:"pedestrian"`
	Footway       string `json:"footway"`
	Cycleway      string `json:"cycleway"`
	Highway       string `json:"highway"`
	Path          string `json:"path"`
	Suburb        string `json:"suburb"`
	City          string `json:"city"`
	Town          string `json:"town"`
	Village       string `json:"village"`
	Hamlet        string `json:"hamlet"`
	County        string `json:"county"`
	Country       string `json:"country"`
	CountryCode   string `json:"country_code"`
	State         string `json:"state"`
	StateDistrict string `json:"state_district"`
	Postcode      string `json:"postcode"`
}

// Locality checks different fields for the locality name
func (a *Address) getCity() string {
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

// Street checks different fields for the street name
func (a *Address) getStreet() string {
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

func (lr *LocationIQResponse) GetAddress() string {
	name := truncateAtFirstComma(lr.DisplayName)
	address := lr.Address.getAddress()
	if idx := strings.Index(address, name); idx == -1 {
		address = fmt.Sprintf("%s, %s", name, address)
	}
	return address
}

func (lr *LocationIQResponse) GetLat() float64 {
	lat, err := strconv.ParseFloat(lr.Lat, 64)
	if err != nil {
		return INVALID_LAT
	}
	return lat
}

func (lr *LocationIQResponse) GetLng() float64 {
	lng, err := strconv.ParseFloat(lr.Lng, 64)
	if err != nil {
		return INVALID_LNG
	}
	return lng
}

func (lr *LocationIQResponse) GetOsmID() int64 {
	osmID, err := strconv.ParseInt(lr.OsmID, 10, 64)
	if err != nil {
		return INVALID_OSM_ID
	}
	return osmID
}

func (lr *LocationIQResponse) GetOsmType() OsmType {
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
