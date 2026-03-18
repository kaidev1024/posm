package posm

import (
	"fmt"
	"strings"
)

type PlaceType uint8

const (
	PlaceTypeNone PlaceType = iota
	PlaceTypePlace
	PlaceTypeOsmNode
	PlaceTypeOsmWay
	PlaceTypeOsmRelation
)

func SanitizeAddress(address string) string {
	if address == "" {
		return address
	}
	address = strings.Join(strings.Fields(address), " ")
	return strings.ReplaceAll(address, " ,", ",")
}

func contructPlaceID(address, lat, lng string) string {
	addressName := strings.ReplaceAll(address, " ", "_")
	return fmt.Sprintf("%s_%s_%s", addressName, lat, lng)
}

func ContructPlaceID(address string, lat, lng float64) string {
	addressName := strings.ReplaceAll(address, " ", "_")
	return fmt.Sprintf("%s_%f_%f", addressName, lat, lng)
}
