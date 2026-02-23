package posm

type OsmCity struct {
	OsmID       int64
	OsmType     OsmType
	Lat         float64
	Lng         float64
	DisplayName string
	Address     string
}

type OsmPoint struct {
	OsmID         int64
	OsmType       OsmType
	Lat           float64
	Lng           float64
	DisplayName   string
	Address       string
	StreetAddress string
	CityAddress   string
}

type OsmStreet struct {
	OsmID       int64
	OsmType     OsmType
	Lat         float64
	Lng         float64
	DisplayName string
	Address     string
}
