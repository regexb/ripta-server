package postgres

import (
	"database/sql"
	"strings"

	"github.com/begizi/ripta-server/stop"
	"github.com/jmoiron/sqlx"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

type dbstop struct {
	ID          string         `db:"stop_id"`
	Description sql.NullString `db:"stop_desc"`
	Name        sql.NullString `db:"stop_name"`
}

type stopStore struct {
	db    *sqlx.DB
	sopts stop.Filter
}

// NewStopStore returns an implimentation of the stop.Store interface
// with mongo as the backing db
func NewStopStore(db *sql.DB) stop.Store {
	return &stopStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (s stopStore) Filter(opts ...stop.FilterOption) stop.Store {
	for _, opt := range opts {
		opt(&s.sopts)
	}
	return &s
}

// List returns a list of all of the stops in the database
func (s *stopStore) List() ([]*stop.Stop, error) {
	query, _, _ := sq.Select("stop_id, stop_name, stop_desc").From("stops").ToSql()

	var stops []*dbstop
	err := s.db.Select(&stops, query)
	if err != nil {
		return nil, err
	}

	var results []*stop.Stop
	for _, s := range stops {
		results = append(results, &stop.Stop{
			ID:          strings.TrimSpace(s.ID),
			Name:        s.Name.String,
			Description: s.Description.String,
		})
	}

	return results, nil
}

// QueryByLocation retuns a list of all stops sorted by distance to
// the provided lat and long
func (s *stopStore) QueryByLocation(lat, long float64) ([]*stop.Stop, error) {
	// hard code distance scope for now
	scope := 200

	queryBuilder := sq.
		Select("DISTINCT stops.stop_id, stops.stop_name, stops.stop_desc").
		From("stops").
		LeftJoin("stop_times ON stop_times.stop_id = stops.stop_id").
		LeftJoin("trips ON trips.trip_id = stop_times.trip_id").
		Where("ST_DWithin(Geography(ST_MakePoint(?,?)), stop_geog, ?)", long, lat, scope).
		PlaceholderFormat(sq.Dollar)

	if s.sopts.Route != "" {
		// queryBuilder = queryBuilder.Where(sq.Eq{"trips.direction_id": "B'?'"}, 0)
		queryBuilder = queryBuilder.Where(sq.Eq{"trips.route_id": s.sopts.Route})
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var stops []*dbstop
	err = s.db.Select(&stops, query, args...)
	if err != nil {
		return nil, err
	}

	var results []*stop.Stop
	for _, s := range stops {
		results = append(results, &stop.Stop{
			ID:          strings.TrimSpace(s.ID),
			Name:        s.Name.String,
			Description: s.Description.String,
		})
	}

	return results, nil
}

// GetByID returns a single stop by the db id. Will return nil
// for unfound stops
func (s *stopStore) GetByID(id string) (*stop.Stop, error) {
	return nil, stop.ErrRecordNotFound
}

// GetByStopID returns a single stop by the stop_id. Will return
// nil for unfound stops
func (s *stopStore) GetByStopID(stopID string) (*stop.Stop, error) {
	return nil, stop.ErrRecordNotFound
}
