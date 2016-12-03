package stop

import (
	"errors"
)

var (
	// ErrRecordNotFound specifies a record does not exist in the store
	ErrRecordNotFound = errors.New("Record not found")
)

// Filter used for query
type Filter struct {
	Route     string
	Direction string
}

// FilterOption type for setting up filters for queries
type FilterOption func(*Filter)

// ContainsRoute filters a query to return only stops that contain a specific route
func ByRoute(route string) FilterOption {
	return func(o *Filter) {
		o.Route = route
	}
}

func ByDirection(direction string) FilterOption {
	return func(o *Filter) {
		o.Direction = direction
	}
}

// Store interface for bus stops
type Store interface {
	List() ([]*Stop, error)
	GetByID(id string) (*Stop, error)
	GetByStopID(stopID string) (*Stop, error)
	QueryByLocation(lat, long float64) ([]*Stop, error)
	Filter(options ...FilterOption) Store
}
