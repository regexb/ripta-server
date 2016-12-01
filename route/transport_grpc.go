package route

import (
	"github.com/begizi/ripta-server/pb"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
)

// GRPCRouteServer defines the grpc implimentation for the route server
type GRPCRouteServer struct {
	stopsByStopID grpctransport.Handler
}

// MakeStopGRPCServer returns a new instance of the grpc server
func MakeRouteGRPCServer(ctx context.Context, endpoints Endpoints, logger log.Logger) *GRPCRouteServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &GRPCRouteServer{
		stopsByStopID: grpctransport.NewServer(
			ctx,
			endpoints.StopsByStopIDEndpoint,
			DecodeGRPCStopsByStopIDRequest,
			EncodeGRPCStopsResponse,
			options...,
		),
	}
}

func (s *GRPCRouteServer) RouteStopsByStopID(ctx context.Context, req *pb.RouteStopRequest) (*pb.RouteStopsResponse, error) {
	_, rep, err := s.stopsByStopID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.RouteStopsResponse), nil
}

// DecodeGRPCStopsByStopIDRequest decodes request from pb type to service type
func DecodeGRPCStopsByStopIDRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.RouteStopRequest)
	return stopByStopIDRequest{req.StopId, req.RouteId}, nil
}

// EncodeGRPCStopsResponse encodes response from service type to pb type
func EncodeGRPCStopsResponse(_ context.Context, response interface{}) (interface{}, error) {
	routeStopResponse := response.([]*Stop)

	pbStops := []*pb.RouteStop{}

	for _, s := range routeStopResponse {
		pbStops = append(pbStops, &pb.RouteStop{
			StopId: s.StopID,
			Time: &pb.RouteStopTime{
				ScheduleTime: s.Time.ScheduleTime.String(),
				ArrivalTime:  s.Time.ScheduleTime.String(),
			},
		})
	}

	return &pb.RouteStopsResponse{pbStops}, nil
}
