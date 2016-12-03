package stop

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

func (mw serviceLoggingMiddleware) List(ctx context.Context) (resp *ListResponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "List",
			"layer", "service",
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.List(ctx)
}

func (mw serviceLoggingMiddleware) QueryByLocation(ctx context.Context, req QueryByLocationRequest) (resp *ListResponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "QueryByLocation",
			"layer", "service",
			"lat", req.Lat,
			"long", req.Long,
			"route", req.Route,
			"direction", req.Direction,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.QueryByLocation(ctx, req)
}

func (mw serviceLoggingMiddleware) Get(ctx context.Context, id string) (resp *GetResponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Get",
			"layer", "service",
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Get(ctx, id)
}

func (mw serviceLoggingMiddleware) GetByStopID(ctx context.Context, stopID string) (resp *GetResponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetByStopID",
			"layer", "service",
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.GetByStopID(ctx, stopID)
}
