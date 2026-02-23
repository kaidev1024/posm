package posm

import (
	"fmt"
	"strconv"
)

// MergeOsmIDAndType combines OSM ID and type into a single string format like "N5200886615" or "W123456"
func CreateOsmTID(osmID int64, osmType OsmType) string {
	var prefix string
	switch osmType {
	case OsmTypeNode:
		prefix = "N"
	case OsmTypeWay:
		prefix = "W"
	case OsmTypeRelation:
		prefix = "R"
	default:
		return ""
	}
	return fmt.Sprintf("%s%d", prefix, osmID)
}

// ParseOsmTID parses a TID string like "N5200886615" into OSM ID and type
func ParseOsmTID(tid string) (int64, OsmType, error) {
	if len(tid) < 2 {
		return 0, OsmTypeNone, fmt.Errorf("invalid tid: too short")
	}

	prefix := tid[0]
	idStr := tid[1:]

	osmID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, OsmTypeNone, fmt.Errorf("invalid osm id in tid: %w", err)
	}

	var osmType OsmType
	switch prefix {
	case 'N', 'n':
		osmType = OsmTypeNode
	case 'W', 'w':
		osmType = OsmTypeWay
	case 'R', 'r':
		osmType = OsmTypeRelation
	default:
		return 0, OsmTypeNone, fmt.Errorf("invalid osm type prefix: %c", prefix)
	}

	return osmID, osmType, nil
}
