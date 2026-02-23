package posm

import (
	"fmt"
	"strings"
)

// address contains address fields specific to OpenStreetMap
type address struct {
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
func (a *address) getCity() string {
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
func (a *address) getStreet() string {
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

func (a *address) isCity() bool {
	if a == nil {
		return false
	}
	return a.getStreet() == "" && a.getCity() != ""
}

func (a *address) getAddress() string {
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
