package geocode

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

// MakeGeocodeHTTPServer returns a new HTTP server for handling geocode requests
func MakeGeocodeHTTPServer(ctx context.Context, endpoints Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorLogger(logger),
	}
	return httptransport.NewServer(
		ctx,
		endpoints.GeocodeEndpoint,
		DecodeHTTPGeocodeRequest,
		EncodeHTTPGeocodeResponse,
		options...,
	)
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

// DecodeHTTPGeocodeRequest processes the HTTP request and returns the service request value
func DecodeHTTPGeocodeRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	query := r.URL.Query().Get("q")

	if query == "" {
		return nil, fmt.Errorf("A search query is required. Got \"%s\"", query)
	}

	req := Request{
		Query: query,
	}

	return req, nil
}

// EncodeHTTPGeocodeResponse processes the service response and writes it to the http response writter
func EncodeHTTPGeocodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
