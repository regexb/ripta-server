package stop

import (
	"golang.org/x/net/context"
)

// Service defines the interface for what a service is for stops
type Service interface {
	List(ctx context.Context) (*ListResponse, error)
	QueryByLocation(ctx context.Context, req QueryByLocationRequest) (*ListResponse, error)
	Get(ctx context.Context, id string) (*GetResponse, error)
	GetByStopID(ctx context.Context, stopID string) (*GetResponse, error)
}

// NewStopService returns a new instance of the Stop service
func NewStopService(store Store) Service {
	return &stopService{
		store: store,
	}
}

type stopService struct {
	store Store
}

func (s stopService) List(_ context.Context) (*ListResponse, error) {
	stops, err := s.store.List()
	if err != nil {
		return nil, err
	}
	return &ListResponse{stops}, nil
}

func (s stopService) QueryByLocation(_ context.Context, req QueryByLocationRequest) (*ListResponse, error) {
	stops, err := s.store.Scope(ContainsRoute(req.Route)).QueryByLocation(req.Lat, req.Long)
	if err != nil {
		return nil, err
	}
	return &ListResponse{stops}, nil
}

func (s stopService) Get(_ context.Context, id string) (*GetResponse, error) {
	stop, err := s.store.GetByID(id)
	if err != nil {
		return nil, err
	}
	return &GetResponse{stop}, nil
}

func (s stopService) GetByStopID(_ context.Context, stopID string) (*GetResponse, error) {
	stop, err := s.store.GetByStopID(stopID)
	if err != nil {
		return nil, err
	}
	return &GetResponse{stop}, nil
}
