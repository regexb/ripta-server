package route

import (
	"time"
)

type Bus struct {
	RouteID        string `json:"route_id"`
	MapsDuration   interface{}
	MapsDistance   interface{}
	RemainingStops []*Stop `json:"remaining_stops"` // Stops left before "destination"
}

// Trip contains the data about the buses location, next stop, time of arrival to stops
type Trip struct {
	RouteID string  `json:"route_id"`
	Stops   []*Stop `json:"stops"`
}

// Stop contains the information about the stop in the sequence of this trip
type Stop struct {
	StopID string    `json:"stop_id"`
	Time   *StopTime `json:"stop_time"`
}

type StopTime struct {
	ScheduleTime time.Time `json:"schedule_time"`
	ArrivalTime  time.Time `json:"arrival_time"`
}
