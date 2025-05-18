package posm

import (
	"fmt"
	"net/http"
)

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

func (a *Address) GetAddress() string {
	return fmt.Sprintf("%s %s, %s, %s %s", a.HouseNumber, a.getStreet(), a.getCity(), a.State, a.Postcode)
}

type LocationIQResponse struct {
	DisplayName string   `json:"display_name"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	Address     *Address `json:"address"`
}
