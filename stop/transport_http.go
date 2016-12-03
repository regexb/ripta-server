package stop

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// MakeStopHTTPServer returns a new HTTP server for handling stop requests
func MakeStopHTTPServer(ctx context.Context, endpoints Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorLogger(logger),
	}
	m := mux.NewRouter()
	m.Handle("/api/v1/stop", httptransport.NewServer(
		ctx,
		endpoints.QueryByLocationEndpoint,
		DecodeHTTPQueryByLocationRequest,
		EncodeHTTPJSONResponse,
		options...,
	)).Queries("location", "{location}")

	m.Handle("/api/v1/stop", httptransport.NewServer(
		ctx,
		endpoints.ListEndpoint,
		DecodeHTTPListRequest,
		EncodeHTTPJSONResponse,
		options...,
	))

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

// DecodeHTTPListRequest processes the HTTP request and returns the service request value
func DecodeHTTPListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return listRequest{}, nil
}

// DecodeHTTPQueryByLocationRequest processes the HTTP request and returns the service request value
func DecodeHTTPQueryByLocationRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	route := r.URL.Query().Get("route")
	direction := r.URL.Query().Get("direction")
	rawLocation := mux.Vars(r)["location"]
	locationArray := strings.Split(rawLocation, ",")

	if len(locationArray) != 2 {
		return nil, fmt.Errorf("location must contain a lat and long (?location=lat,long). Got \"%s\"", rawLocation)
	}

	// Parse Lat
	lat, err := strconv.ParseFloat(locationArray[0], 64)
	if err != nil {
		return nil, fmt.Errorf("lat must be a float. Got \"%s\"", locationArray[0])
	}

	// Parse Long
	long, err := strconv.ParseFloat(locationArray[1], 64)
	if err != nil {
		return nil, fmt.Errorf("long must be a float. Got \"%s\"", locationArray[1])
	}

	return QueryByLocationRequest{route, lat, long, direction}, nil
}

// EncodeHTTPJSONResponse processes the service response and writes it to the http response writter
func EncodeHTTPJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
