package geocode

import (
	"github.com/begizi/ripta-server/pb"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
)

// GRPCGeocodeServer defines the grpc implimentation for the stop server
type GRPCGeocodeServer struct {
	geocode grpctransport.Handler
}

// MakeGeocodeGRPCServer returns a new instance of the grpc server
func MakeGeocodeGRPCServer(ctx context.Context, endpoints Endpoints, logger log.Logger) *GRPCGeocodeServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &GRPCGeocodeServer{
		geocode: grpctransport.NewServer(
			ctx,
			endpoints.GeocodeEndpoint,
			DecodeGRPCGeocodeRequest,
			EncodeGRPCGeocodeResponse,
			options...,
		),
	}
}

// Geocode transport handler
func (s *GRPCGeocodeServer) Geocode(ctx context.Context, req *pb.GeocodeRequest) (*pb.GeocodeResponse, error) {
	_, rep, err := s.geocode.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GeocodeResponse), nil
}

// DecodeGRPCGeocodeRequest decodes request from pb type to service type
func DecodeGRPCGeocodeRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GeocodeRequest)
	return Request{req.Query}, nil
}

// EncodeGRPCGeocodeResponse encodes response from service type to pb type
func EncodeGRPCGeocodeResponse(_ context.Context, response interface{}) (interface{}, error) {
	geocodeResponse := response.(*Response)
	pbLocations := []*pb.Geocode{}

	for _, g := range geocodeResponse.Data {
		pbLocations = append(pbLocations, &pb.Geocode{
			Address: g.Address,
			Lat:     g.Lat,
			Long:    g.Long,
		})
	}
	return &pb.GeocodeResponse{pbLocations}, nil
}
