package realtime

import (
	"strconv"
	"time"
)

type TimeStamp time.Time

func (t *TimeStamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	*t = TimeStamp(time.Unix(int64(ts), 0))

	return nil
}

type TripUpdates struct {
	Header     *Header        `json:"header"`
	TripEntity []*TripEntitiy `json:"entity"`
}

type Header struct {
	Timestamp TimeStamp `json:"timestamp"`
}

type TripEntitiy struct {
	ID         string      `json:"id"`
	IsDeleted  bool        `json:"is_deleted"`
	Vehicle    *Vehicle    `json:"vehicle"`
	Alert      *Alert      `json:"alert"`
	TripUpdate *TripUpdate `json:"trip_update"`
}

type TripUpdate struct {
	Trip            *Trip             `json:"trip"`
	StopTimeUpdates []*StopTimeUpdate `json:"stop_time_update"`
	Vehicle         *Vehicle          `json:"vehicle"`
	TimeStamp       TimeStamp         `json:"timestamp"`
}

type Trip struct {
	TripID               string `json:"trip_id"`
	StartTime            string `json:"start_time"`
	StartDate            string `json:"start_date"`
	ScheduleRelationship int    `json:"schedule_relationship"`
	RouteID              string `json:"route_id"`
}

type StopTime struct {
	Delay       int       `json:"delay"`
	Time        TimeStamp `json:"time"`
	Uncertainty int       `json:"uncertainty"`
}

type StopTimeUpdate struct {
	StopID               string    `json:"stop_id"`
	StopSequence         int       `json:"stop_sequence"`
	ScheduleRelationship int       `json:"schedule_relationship"`
	Arrival              *StopTime `json:"arrival"`
	Departure            *StopTime `json:"departure"`
}

type TripUpdatesService service

func (s *TripUpdatesService) Get() (*TripUpdates, error) {
	req, err := s.client.NewRequest("GET", RIPTAAPITripUpdates, nil)
	if err != nil {
		return nil, err
	}

	var tripResp = new(TripUpdates)
	if _, err := s.client.Do(req, tripResp); err != nil {
		return nil, err
	}

	return tripResp, nil
}
