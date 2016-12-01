package geocode

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

// Middleware definition for wrapping the service
type Middleware func(Service) Service

// EndpointLoggingMiddleware middleware for logging endpoint response times
func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log("layer", "endpoint", "error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// ServiceLoggingMiddleware middleware for logging service time and response data
func ServiceLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return serviceLoggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type serviceLoggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw serviceLoggingMiddleware) Geocode(ctx context.Context, req Request) (resp *Response, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Geocode",
			"layer", "service",
			"geo_query", req.Query,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Geocode(ctx, req)
}
