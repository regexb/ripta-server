package route

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

// Endpoints for interfacing with the service
type Endpoints struct {
	StopsByStopIDEndpoint endpoint.Endpoint
}

type stopByStopIDRequest struct {
	stopID  string
	routeID string
}

// StopsByStopID  StopsByStopIDEndpoint
func (e Endpoints) StopsByStopID(ctx context.Context, request stopByStopIDRequest) ([]*Stop, error) {
	resp, err := e.StopsByStopIDEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.([]*Stop), nil
}

// MakeStopsByStopIDEndpoint returns an implimentation of the endpoint
func MakeStopsByStopIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(stopByStopIDRequest)
		resp, err := s.StopsByStopID(ctx, req.routeID, req.stopID)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}
