package stop

import (
	"errors"
)

var (
	// ErrRecordNotFound specifies a record does not exist in the store
	ErrRecordNotFound = errors.New("Record not found")
)

// Scope used for query
type Scope struct {
	Route string
}

// ScopeOption type for setting up scopes for queries
type ScopeOption func(*Scope)

// ContainsRoute scopes a query to return only stops that contain a specific route
func ContainsRoute(route string) ScopeOption {
	return func(o *Scope) {
		o.Route = route
	}
}

// Store interface for bus stops
type Store interface {
	List() ([]*Stop, error)
	GetByID(id string) (*Stop, error)
	GetByStopID(stopID string) (*Stop, error)
	QueryByLocation(lat, long float64) ([]*Stop, error)
	Scope(options ...ScopeOption) Store
}
