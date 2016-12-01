package stop

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

// Endpoints for interfacing with the service
type Endpoints struct {
	ListEndpoint            endpoint.Endpoint
	QueryByLocationEndpoint endpoint.Endpoint
}

type listRequest struct{}

// List  ListEndpoint
func (e Endpoints) List(ctx context.Context) (*ListResponse, error) {
	resp, err := e.ListEndpoint(ctx, listRequest{})
	if err != nil {
		return nil, err
	}
	return resp.(*ListResponse), nil
}

// MakeListEndpoint returns an implimentation of the endpoint
func MakeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (response interface{}, err error) {
		resp, err := s.List(ctx)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

// QueryByLocation  ListEndpoint
func (e Endpoints) QueryByLocation(ctx context.Context, r QueryByLocationRequest) (*ListResponse, error) {
	resp, err := e.QueryByLocationEndpoint(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*ListResponse), nil
}

// MakeQueryByLocationEndpoint returns an implimentation of the endpoint
func MakeQueryByLocationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		queryReq := request.(QueryByLocationRequest)
		resp, err := s.QueryByLocation(ctx, queryReq)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}
