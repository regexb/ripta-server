package realtime

// Filter used for query
type Filter struct {
	Route string
	Stop  string
}

// FilterOption type for setting up filters for queries
type FilterOption func(*Filter)

// ContainsRoute filters a query to return only stops that contain a specific route
func ByRoute(route string) FilterOption {
	return func(o *Filter) {
		o.Route = route
	}
}

func ByStop(stopID string) FilterOption {
	return func(o *Filter) {
		o.Stop = stopID
	}
}

type Store interface {
	GetTripUpdates() ([]*TripUpdate, error)
	Filter(options ...FilterOption) Store
}
