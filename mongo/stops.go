package mongo

import (
	"github.com/begizi/ripta-server/stop"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	stopDoc = "stop"
)

type geoJSON struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

type dbstop struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	StopID   string        `json:"stop_id"`
	Name     string        `json:"stop_name"`
	Location geoJSON       `json:"location"`
	Routes   []string      `json:"route_ids"`
}

type stopStore struct {
	db      string
	session *mgo.Session
	sopts   stop.Filter
}

// NewStopStore returns an implimentation of the stop.Store interface
// with mongo as the backing db
func NewStopStore(db string, session *mgo.Session) (stop.Store, error) {
	s := &stopStore{
		db:      db,
		session: session,
	}

	stopIndex := mgo.Index{
		Key:        []string{"stop_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	geoIndex := mgo.Index{
		Key: []string{"$2dsphere:location"},
	}

	sess := s.session.Copy()
	defer sess.Close()

	c := sess.DB(s.db).C(stopDoc)

	if err := c.EnsureIndex(stopIndex); err != nil {
		return nil, err
	}

	if err := c.EnsureIndex(geoIndex); err != nil {
		return nil, err
	}

	return s, nil
}

func (s stopStore) Filter(opts ...stop.FilterOption) stop.Store {
	for _, opt := range opts {
		opt(&s.sopts)
	}
	return &s
}

// List returns a list of all of the stops in the database
func (s *stopStore) List() ([]*stop.Stop, error) {
	sess := s.session.Copy()
	defer sess.Close()

	query := bson.M{}

	if s.sopts.Route != "" {
		query = bson.M{
			"routes": s.sopts.Route,
		}
	}

	var stops []*dbstop
	c := sess.DB(s.db).C(stopDoc)
	err := c.Find(query).All(&stops)
	if err != nil {
		return nil, stop.ErrRecordNotFound
	}

	var results []*stop.Stop
	for _, s := range stops {
		results = append(results, &stop.Stop{
			ID:   s.StopID,
			Name: s.Name,
		})
	}

	return results, nil
}

// QueryByLocation retuns a list of all stops sorted by distance to
// the provided lat and long
func (s *stopStore) QueryByLocation(lat, long float64) ([]*stop.Stop, error) {
	sess := s.session.Copy()
	defer sess.Close()

	// hard code distance scope for now
	scope := 150

	query := bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
				"$maxDistance": scope,
			},
		},
	}

	if s.sopts.Route != "" {
		query["routes"] = s.sopts.Route
	}

	var stops []*dbstop
	c := sess.DB(s.db).C(stopDoc)
	err := c.Find(query).Limit(5).All(&stops)
	if err != nil {
		return nil, err
	}

	var results []*stop.Stop
	for _, s := range stops {
		results = append(results, &stop.Stop{
			ID:   s.StopID,
			Name: s.Name,
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
