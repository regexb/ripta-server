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
}

func NewRealtimeStore(client *realtime.Client, exp time.Duration) realtime.Store {
	return &InmemRealtimeStore{
		client:  client,
		expTime: exp,
	}
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

func (r *InmemRealtimeStore) GetTripUpdates() ([]*realtime.TripUpdate, error) {
	tripupdates, err := r.getTripUpdates()
	if err != nil {
		return nil, err
	}

	updates := []*realtime.TripUpdate{}
	for _, t := range tripupdates.TripEntity {
		updates = append(updates, t.TripUpdate)
	}
	return updates, nil
}

func (r *InmemRealtimeStore) GetTripUpdatesByRouteID(routeID string) ([]*realtime.TripUpdate, error) {
	tripupdates, err := r.getTripUpdates()
	if err != nil {
		return nil, err
	}

	updates := []*realtime.TripUpdate{}
	for _, t := range tripupdates.TripEntity {
		if t.TripUpdate.Trip.RouteID == routeID {
			updates = append(updates, t.TripUpdate)
		}
	}
	return updates, nil
}
