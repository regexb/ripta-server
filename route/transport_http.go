package route

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// MakeRouteHTTPServer returns a new HTTP server for handling geocode requests
func MakeRouteHTTPServer(ctx context.Context, endpoints Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorLogger(logger),
	}
	m := mux.NewRouter()
	m.Handle("/api/v1/route/{routeID}", httptransport.NewServer(
		ctx,
		endpoints.StopsByStopIDEndpoint,
		DecodeHTTPStopsByStopIDRequest,
		EncodeHTTPJSONResponse,
		options...,
	)).Queries("stop_id", "{stopID}")
	return m
}

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	msg := err.Error()

	if e, ok := err.(httptransport.Error); ok {
		msg = e.Err.Error()
		switch e.Domain {
		case httptransport.DomainDecode:
			code = http.StatusBadRequest

		case httptransport.DomainDo:
			code = http.StatusBadRequest
		}
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorWrapper{Error: msg})
}

// DecodeHTTPStopsByStopIDRequest processes the HTTP request and returns the service request value
func DecodeHTTPStopsByStopIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	stopID := mux.Vars(r)["stopID"]
	routeID := mux.Vars(r)["routeID"]

	req := stopByStopIDRequest{
		stopID:  stopID,
		routeID: routeID,
	}

	return req, nil
}

// EncodeHTTPJSONResponse processes the service response and writes it to the http response writter
func EncodeHTTPJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
