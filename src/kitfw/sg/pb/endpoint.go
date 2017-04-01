package pb

// This file contains methods to make individual endpoints from services,
// request and response types to serve those endpoints, as well as encoders and
// decoders for those types, for all of our supported transport serialization
// formats. It also includes endpoint middlewares.

import (
	"fmt"
	kitfwService "kitfw/sg/service"
	"time"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type Endpoint struct {
	process endpoint.Endpoint
}

// Sum implements Service. Primarily useful in a client.
func (e Endpoint) Process(ctx context.Context, request *KitfwRequest) (*KitfwReply, error) {
	response, err := e.process(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(*KitfwReply), err
}

func MakeProcessEndpoint(s kitfwService.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		q := request.(*KitfwRequest)
		ctx = context.WithValue(ctx, "logid", q.Logid)
		v, err := s.Process(ctx, q.Protoid, q.Payload)
		if err != nil {
			return nil, err
		}
		return &KitfwReply{
			Protoid: q.Protoid + 1,
			Logid:   q.Logid,
			Payload: v,
		}, nil
	}
}

func EndpointInstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)

		}
	}
}

func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log("error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)

		}
	}
}
