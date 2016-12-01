package geocode

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

// Endpoints for interfacing with the service
type Endpoints struct {
	GeocodeEndpoint endpoint.Endpoint
}

// Geocode  GeocodeEndpoint
func (e Endpoints) Geocode(ctx context.Context, request Request) (*Response, error) {
	resp, err := e.GeocodeEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*Response), nil
}

// MakeGeocodeEndpoint returns an implimentation of the endpoint
func MakeGeocodeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Request)
		resp, err := s.Geocode(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}
