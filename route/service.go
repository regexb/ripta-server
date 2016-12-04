package route

import (
	"github.com/begizi/ripta-server/realtime"
	"golang.org/x/net/context"
	"time"
)

type Service interface {
	StopsByStopID(ctx context.Context, routeID, stopID string) ([]*Stop, error) // Get the stops for a route and a specific stop
	// GetBusByRouteAndStop(routeID, stopID string) (*Bus, error)      // Get the next bus for a route and a specific stop
	// GetStopsByRoute(routeID string) ([]*Stop, error)                // Get all the stops for a specific route
	// GetTripsbyRoute(routeID string) ([]*Trip, error)                // Get all of the next busses to a specific route
	// GetStopsByStop(stopID string) ([]*Stop, error)                  // Get all of the stops for a specific stop
	// GetTripsByStop(stopID string) ([]*Trip, error)                  // Get all of the next buses to a specific stop
}

func NewRouteService(store realtime.Store) Service {
	return &routeService{
		store: store,
	}
}

type routeService struct {
	store realtime.Store
}

func (s *routeService) StopsByStopID(ctx context.Context, routeID, stopID string) ([]*Stop, error) {
	tripUpdates, err := s.store.Filter(realtime.ByRoute(routeID), realtime.ByStop(stopID)).GetTripUpdates()
	if err != nil {
		return nil, err
	}

	stops := []*Stop{}
	for _, tripUpdate := range tripUpdates {
		for _, stopTimeUpdate := range tripUpdate.StopTimeUpdates {
			if stopTimeUpdate.Arrival != nil {
				arrivalTime := time.Time(stopTimeUpdate.Arrival.Time)
				timeWithDelay := arrivalTime.Add(time.Duration(stopTimeUpdate.Arrival.Delay) * time.Second)
				stops = append(stops, &Stop{
					StopID: stopTimeUpdate.StopID,
					Time: &StopTime{
						ScheduleTime: arrivalTime,
						ArrivalTime:  timeWithDelay,
					},
				})
			}
		}
	}

	return stops, nil
}
