package stop

import (
	"github.com/begizi/ripta-server/pb"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
)

// GRPCServer defines the grpc implimentation for the stop server
type GRPCStopServer struct {
	list          grpctransport.Handler
	getByLocation grpctransport.Handler
}

// MakeStopGRPCServer returns a new instance of the grpc server
func MakeStopGRPCServer(ctx context.Context, endpoints Endpoints, logger log.Logger) *GRPCStopServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &GRPCStopServer{
		list: grpctransport.NewServer(
			ctx,
			endpoints.ListEndpoint,
			DecodeGRPCListRequest,
			EncodeGRPCListResponse,
			options...,
		),
		getByLocation: grpctransport.NewServer(
			ctx,
			endpoints.QueryByLocationEndpoint,
			DecodeGRPCGetByLocationRequest,
			EncodeGRPCListResponse,
			options...,
		),
	}
}

func (s *GRPCStopServer) GetStopsByLocation(ctx context.Context, req *pb.StopLocationRequest) (*pb.StopsResponse, error) {
	_, rep, err := s.getByLocation.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.StopsResponse), nil
}

// DecodeGRPCListRequest decodes request from pb type to service type
func DecodeGRPCGetByLocationRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.StopLocationRequest)
	return QueryByLocationRequest{req.Route, req.Lat, req.Long}, nil
}

// List transport handler
func (s *GRPCStopServer) ListStops(ctx context.Context, req *pb.ListStopsRequest) (*pb.StopsResponse, error) {
	_, rep, err := s.list.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.StopsResponse), nil
}

// DecodeGRPCListRequest decodes request from pb type to service type
func DecodeGRPCListRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return listRequest{}, nil
}

// EncodeGRPCListResponse encodes response from service type to pb type
func EncodeGRPCListResponse(_ context.Context, response interface{}) (interface{}, error) {
	listResponse := response.(*ListResponse)
	pbStops := []*pb.Stop{}

	for _, s := range listResponse.Data {
		pbStops = append(pbStops, &pb.Stop{
			Id:     s.ID.String(),
			StopId: s.StopID,
			Name:   s.Name,
			Location: &pb.Location{
				Lat:  s.Location.Coordinates[0],
				Long: s.Location.Coordinates[1],
			},
			Routes: s.Routes,
		})
	}
	return &pb.StopsResponse{pbStops}, nil
}
