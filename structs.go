package posm

type OsmCity struct {
	PlaceID     string
	Lat         float64
	Lng         float64
	DisplayName string
	Address     string
}

type OsmPoint struct {
	PlaceID          string
	Lat              float64
	Lng              float64
	DisplayName      string
	Address          string
	StreetSearchText string
	CitySearchText   string
}

type OsmStreet struct {
	PlaceID     string
	Lat         float64
	Lng         float64
	DisplayName string
	Address     string
}
