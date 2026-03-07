package posm

import (
	"fmt"
	"strconv"
	"strings"
)

type locationIQResponse struct {
	PlaceID     string   `json:"place_id"`
	OsmID       string   `json:"osm_id"`
	OsmType     string   `json:"osm_type"`
	DisplayName string   `json:"display_name"`
	Lat         string   `json:"lat"`
	Lng         string   `json:"lon"`
	Address     *address `json:"address"`
}

func (lr *locationIQResponse) getPointAddress() string {
	if lr == nil {
		return ""
	}
	if lr.Address == nil {
		return lr.DisplayName
	}
	return lr.Address.getAddress()
}

func (lr *locationIQResponse) getCityAddress() string {
	if lr == nil {
		return ""
	}
	if lr.Address == nil {
		return lr.DisplayName
	}
	address := lr.Address
	city := address.getCity()
	if city == "" {
		city = address.County
	}
	return fmt.Sprintf("%s, %s", city, address.State)
}

func (lr *locationIQResponse) getStreetAddress() string {
	if lr == nil {
		return ""
	}
	if lr.Address == nil {
		return lr.DisplayName
	}
	address := lr.Address
	return fmt.Sprintf("%s, %s, %s", address.getStreet(), address.getCity(), address.State)
}

func (lr *locationIQResponse) getStreetSearchText() string {
	if lr == nil {
		return ""
	}
	if lr.Address == nil {
		return lr.DisplayName
	}
	address := lr.Address
	street := address.getStreet()
	if street == "" {
		return ""
	}
	return fmt.Sprintf("%s, %s, %s, %s", street, address.getCity(), address.State, address.CountryCode)
}

func (lr *locationIQResponse) getCitySearchText() string {
	if lr == nil {
		return ""
	}
	if lr.Address == nil {
		return lr.DisplayName
	}
	address := lr.Address
	city := address.getCity()
	if city == "" {
		return ""
	}
	return fmt.Sprintf("%s, %s, %s", city, address.State, address.CountryCode)
}

func (lr *locationIQResponse) parseCoordinates() (float64, float64, error) {
	if lr == nil {
		return HEADQUARTER_LAT, HEADQUARTER_LNG, fmt.Errorf("empty location")
	}
	lat, err := strconv.ParseFloat(lr.Lat, 64)
	if err != nil {
		return HEADQUARTER_LAT, HEADQUARTER_LNG, err
	}
	lng, err := strconv.ParseFloat(lr.Lng, 64)
	if err != nil {
		return HEADQUARTER_LAT, HEADQUARTER_LNG, err
	}
	return lat, lng, nil
}

func (lr *locationIQResponse) isCity() bool {
	if lr == nil {
		return false
	}
	return lr.Address.isCity()
}

func (lr *locationIQResponse) getPlaceID() string {
	if lr == nil {
		return ""
	}
	prefix := ""
	id := ""
	if lr.OsmType == "" && lr.OsmID == "" {
		id = lr.PlaceID
		if lr.OsmType == "node" {
			prefix = "N"
		}
		if lr.OsmType == "way" {
			prefix = "W"
		}
		if lr.OsmType == "relation" {
			prefix = "R"
		}
	} else if lr.PlaceID != "" {
		id = lr.PlaceID
		prefix = "P"
	} else {
		prefix = "X"
		addressName := ""
		if lr.DisplayName != "" {
			addressName = lr.DisplayName
		} else if lr.Address != nil {
			addressName = lr.Address.getAddress()
		}
		addressName = strings.ReplaceAll(addressName, " ", "_")
		id = fmt.Sprintf("%s_%s_%s", addressName, lr.Lat, lr.Lng)
	}
	return fmt.Sprintf("%s%s", prefix, id)
}
