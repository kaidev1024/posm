package posm

import (
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
