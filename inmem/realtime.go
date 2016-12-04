package inmem

import (
	"github.com/begizi/ripta-server/realtime"
	"time"
)

type InmemRealtimeStore struct {
	expTime        time.Duration
	client         *realtime.Client
	tripupdates    *realtime.TripUpdates
	tripupdatesExp time.Time
	ropts          realtime.Filter
}

func NewRealtimeStore(client *realtime.Client, exp time.Duration) realtime.Store {
	return &InmemRealtimeStore{
		client:  client,
		expTime: exp,
	}
}

func (r InmemRealtimeStore) Filter(opts ...realtime.FilterOption) realtime.Store {
	for _, opt := range opts {
		opt(&r.ropts)
	}
	return &r
}

func (r *InmemRealtimeStore) getTripUpdates() (*realtime.TripUpdates, error) {
	if time.Since(r.tripupdatesExp) > r.expTime {
		t, err := r.client.TripUpdates.Get()
		if err != nil {
			return nil, err
		}
		r.tripupdates = t
		r.tripupdatesExp = time.Now()
	}

	return r.tripupdates, nil
}

func filterTripUpdatesByRoute(tripUpdates []*realtime.TripUpdate, routeID string) []*realtime.TripUpdate {
	updates := []*realtime.TripUpdate{}
	for _, t := range tripUpdates {
		if t.Trip.RouteID == routeID {
			updates = append(updates, t)
		}
	}
	return updates
}

func filterTripUpdatesByStop(tripUpdates []*realtime.TripUpdate, stopID string) []*realtime.TripUpdate {
	updates := []*realtime.TripUpdate{}
	for _, t := range tripUpdates {
		stops := []*realtime.StopTimeUpdate{}
		for _, stopTimeUpdate := range t.StopTimeUpdates {
			if stopTimeUpdate.StopID == stopID {
				stops = append(stops, stopTimeUpdate)
			}
		}
		t.StopTimeUpdates = stops
		updates = append(updates, t)
	}
	return updates
}

func (r *InmemRealtimeStore) GetTripUpdates() ([]*realtime.TripUpdate, error) {
	tripupdates, err := r.getTripUpdates()
	if err != nil {
		return nil, err
	}

	updates := []*realtime.TripUpdate{}
	for _, t := range tripupdates.TripEntity {
		updates = append(updates, t.TripUpdate)
	}

	if r.ropts.Route != "" {
		updates = filterTripUpdatesByRoute(updates, r.ropts.Route)
	}

	if r.ropts.Stop != "" {
		updates = filterTripUpdatesByStop(updates, r.ropts.Stop)
	}

	return updates, nil
}
