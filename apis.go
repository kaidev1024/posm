package posm

import "net/http"

const (
	OsmTypeNone OsmType = iota
	OsmTypeNode
	OsmTypeWay
	OsmTypeRelation
)

const INVALID_LAT float64 = 200.0
const INVALID_LNG float64 = 200.0
const INVALID_OSM_ID int64 = 0
const HEADQUARTER_LAT float64 = 37.7955
const HEADQUARTER_LNG float64 = -122.3937

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new LocationIQ client
func Init(accessToken string) {
	locationIQAccessToken = accessToken
	searchClient = &Client{
		BaseURL:    "https://us1.locationiq.com/v1/search",
		HTTPClient: &http.Client{},
	}
	autoCompleteClient = &Client{
		BaseURL:    "https://api.locationiq.com/v1/autocomplete",
		HTTPClient: &http.Client{},
	}
	lookupClient = &Client{
		BaseURL:    "https://us1.locationiq.com/v1/lookup",
		HTTPClient: &http.Client{},
	}
}
