package posm

import (
	"fmt"
	"strconv"
	"strings"
)

type locationIQResponse struct {
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
	name := truncateAtFirstComma(lr.DisplayName)
	address := lr.Address.getAddress()
	if idx := strings.Index(address, name); idx == -1 {
		address = fmt.Sprintf("%s, %s", name, address)
	}
	return address
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

func (lr *locationIQResponse) getLat() (float64, error) {
	if lr == nil {
		return INVALID_LAT, fmt.Errorf("empty location")
	}
	lat, err := strconv.ParseFloat(lr.Lat, 64)
	if err != nil {
		return INVALID_LAT, err
	}
	return lat, nil
}

func (lr *locationIQResponse) getLng() (float64, error) {
	if lr == nil {
		return INVALID_LNG, fmt.Errorf("empty location")
	}
	lng, err := strconv.ParseFloat(lr.Lng, 64)
	if err != nil {
		return INVALID_LNG, err
	}
	return lng, nil
}

func (lr *locationIQResponse) getOsmID() (int64, error) {
	if lr == nil {
		return INVALID_OSM_ID, fmt.Errorf("empty location")
	}
	if lr.OsmID == "" {
		return INVALID_OSM_ID, nil
	}
	osmID, err := strconv.ParseInt(lr.OsmID, 10, 64)
	if err != nil {
		return INVALID_OSM_ID, err
	}
	return osmID, nil
}

func (lr *locationIQResponse) getOsmType() OsmType {
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

func (lr *locationIQResponse) isCity() bool {
	if lr == nil {
		return false
	}
	return lr.Address.isCity()
}
