package geocode

type Request struct {
	Query string `json:"query"`
}

type Response struct {
	Data []*GeocodeAddress `json:"data"`
}

type GeocodeAddress struct {
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
}
