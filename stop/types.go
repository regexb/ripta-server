package stop

import (
	"gopkg.in/mgo.v2/bson"
)

// ListResponse for returning multiple stops from an endpoint
type ListResponse struct {
	Data []*Stop `json:"data"`
}

type QueryByLocationRequest struct {
	Route string  `json:"route"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

// GetResponse for returning a sindle stop from an endpoint
type GetResponse struct {
	Data *Stop `json:"data"`
}

// GeoJSON contains location coordinates. Coordinates field contains the lat and long pairs
type GeoJSON struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

// Stop contains the bus stop information
type Stop struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	StopID   string        `json:"stop_id"`
	Name     string        `json:"stop_name"`
	Location GeoJSON       `json:"location"`
	Routes   []string      `json:"route_ids"`
}
