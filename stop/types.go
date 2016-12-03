package stop

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

// Stop contains the bus stop information
type Stop struct {
	ID          string `json:"stop_id"`
	Name        string `json:"stop_name"`
	Description string `json:"stop_desc"`
}
