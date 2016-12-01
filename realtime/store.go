package realtime

type Store interface {
	GetTripUpdates() ([]*TripUpdate, error)
	GetTripUpdatesByRouteID(routeID string) ([]*TripUpdate, error)
}
