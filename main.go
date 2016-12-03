package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"googlemaps.github.io/maps"

	"github.com/begizi/ripta-server/geocode"
	"github.com/begizi/ripta-server/inmem"
	"github.com/begizi/ripta-server/pb"
	"github.com/begizi/ripta-server/postgres"
	"github.com/begizi/ripta-server/realtime"
	"github.com/begizi/ripta-server/route"
	"github.com/begizi/ripta-server/stop"
)

// RiptaGRPCServer aggregates all of the grpc servers
type RiptaGRPCServer struct {
	*stop.GRPCStopServer
	*geocode.GRPCGeocodeServer
	*route.GRPCRouteServer
}

func main() {
	port := os.Getenv("HTTP_PORT")
	// default for http port
	if port == "" {
		port = "8080"
	}

	gRPCPort := os.Getenv("GRPC_PORT")
	// default for grpc port
	if gRPCPort == "" {
		gRPCPort = "9001"
	}

	dbAddr := os.Getenv("DB_ADDR")
	// default for db
	if dbAddr == "" {
		dbAddr = "postgresql://postgres@127.0.0.1/ripta?sslmode=disable"
	}

	// Map Client
	mapsAPIKey := os.Getenv("MAPS_API_KEY")
	client, err := maps.NewClient(maps.WithAPIKey(mapsAPIKey))
	if err != nil || client == nil {
		panic(fmt.Errorf("Filed to initialize maps client: Client %v: Error %v", client, err))
	}

	// Setup RealtimeClient
	realtimeClient := realtime.NewClient(nil)

	// Setup DB connection
	db, err := sql.Open("postgres", dbAddr)
	if err != nil {
		panic(fmt.Errorf("Failed to initialize database connection: Error %v", err))
	}
	defer db.Close()

	// Stores
	stopStore := postgres.NewStopStore(db)
	realtimeStore := inmem.NewRealtimeStore(realtimeClient, 60*time.Second)

	// Context
	ctx := context.Background()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	// Business domain.
	var geocodeService geocode.Service
	{
		geocodeService = geocode.NewGeocodeService(*client)
		geocodeService = geocode.ServiceLoggingMiddleware(logger)(geocodeService)
	}

	var stopService stop.Service
	{
		stopService = stop.NewStopService(stopStore)
		stopService = stop.ServiceLoggingMiddleware(logger)(stopService)
	}

	var routeService route.Service
	{
		routeService = route.NewRouteService(realtimeStore)
		routeService = route.ServiceLoggingMiddleware(logger)(routeService)
	}

	// Endpoint domain.
	var geocodeEndpoint endpoint.Endpoint
	{
		geocodeLogger := log.NewContext(logger).With("method", "Geocode")
		geocodeEndpoint = geocode.MakeGeocodeEndpoint(geocodeService)
		geocodeEndpoint = geocode.EndpointLoggingMiddleware(geocodeLogger)(geocodeEndpoint)
	}

	var listEndpoint endpoint.Endpoint
	{
		listLogger := log.NewContext(logger).With("method", "StopList")
		listEndpoint = stop.MakeListEndpoint(stopService)
		listEndpoint = stop.EndpointLoggingMiddleware(listLogger)(listEndpoint)
	}

	var queryByLocationEndpoint endpoint.Endpoint
	{
		queryByLocationLogger := log.NewContext(logger).With("method", "StopQueryByLocation")
		queryByLocationEndpoint = stop.MakeQueryByLocationEndpoint(stopService)
		queryByLocationEndpoint = stop.EndpointLoggingMiddleware(queryByLocationLogger)(queryByLocationEndpoint)
	}

	var stopsByStopIDEndpoint endpoint.Endpoint
	{
		stopsByStopIDLogger := log.NewContext(logger).With("method", "RouteStopsByStopID")
		stopsByStopIDEndpoint = route.MakeStopsByStopIDEndpoint(routeService)
		stopsByStopIDEndpoint = route.EndpointLoggingMiddleware(stopsByStopIDLogger)(stopsByStopIDEndpoint)
	}

	// Error chan
	errc := make(chan error)

	// Interrupt handler
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// HTTP transport
	go func() {
		// HTTP Handler Domain
		var geocodeHandler http.Handler
		var stopHandler http.Handler
		var routeHandler http.Handler
		{
			geocodeEndpoints := geocode.Endpoints{
				GeocodeEndpoint: geocodeEndpoint,
			}
			stopEndpoints := stop.Endpoints{
				ListEndpoint:            listEndpoint,
				QueryByLocationEndpoint: queryByLocationEndpoint,
			}
			routeEndpoints := route.Endpoints{
				StopsByStopIDEndpoint: stopsByStopIDEndpoint,
			}
			logger := log.NewContext(logger).With("transport", "HTTP")
			geocodeHandler = geocode.MakeGeocodeHTTPServer(ctx, geocodeEndpoints, logger)
			stopHandler = stop.MakeStopHTTPServer(ctx, stopEndpoints, logger)
			routeHandler = route.MakeRouteHTTPServer(ctx, routeEndpoints, logger)
		}

		// Root Handler
		r := mux.NewRouter()
		r.Handle("/api/v1/geocode", geocodeHandler)
		r.Handle("/api/v1/stop", stopHandler)
		r.Handle("/api/v1/route/{routeID}", routeHandler)

		logger.Log("msg", "HTTP Server Started", "port", port)
		errc <- http.ListenAndServe(":"+port, r)
	}()

	// gRPC transport
	go func() {
		lis, err := net.Listen("tcp", ":"+gRPCPort)
		if err != nil {
			errc <- err
			return
		}
		defer lis.Close()

		s := grpc.NewServer()

		// Mechanical domain.
		var riptaServer pb.RiptaServer
		{
			stopEndpoints := stop.Endpoints{
				ListEndpoint:            listEndpoint,
				QueryByLocationEndpoint: queryByLocationEndpoint,
			}
			geocodeEndpoints := geocode.Endpoints{
				GeocodeEndpoint: geocodeEndpoint,
			}
			routeEndpoints := route.Endpoints{
				StopsByStopIDEndpoint: stopsByStopIDEndpoint,
			}
			logger := log.NewContext(logger).With("transport", "gRPC")
			stopServer := stop.MakeStopGRPCServer(ctx, stopEndpoints, logger)
			geocodeServer := geocode.MakeGeocodeGRPCServer(ctx, geocodeEndpoints, logger)
			routeServer := route.MakeRouteGRPCServer(ctx, routeEndpoints, logger)

			riptaServer = RiptaGRPCServer{
				GRPCStopServer:    stopServer,
				GRPCGeocodeServer: geocodeServer,
				GRPCRouteServer:   routeServer,
			}

		}

		pb.RegisterRiptaServer(s, riptaServer)

		logger.Log("msg", "GRPC Server Started", "port", gRPCPort)
		errc <- s.Serve(lis)
	}()

	logger.Log("exit", <-errc)
}
