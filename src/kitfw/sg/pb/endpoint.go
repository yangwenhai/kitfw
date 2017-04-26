package pb

// This file contains methods to make individual endpoints from services,
// request and response types to serve those endpoints, as well as encoders and
// decoders for those types, for all of our supported transport serialization
// formats. It also includes endpoint middlewares.

import (
	"errors"
	"fmt"
	logger "kitfw/sg/log"
	protocol "kitfw/sg/protocol"
	kitfwService "kitfw/sg/service"
	"time"

	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	stdopentracing "github.com/opentracing/opentracing-go"
	otext "github.com/opentracing/opentracing-go/ext"
)

func MakeProcessEndpoint(s kitfwService.Service, tracer stdopentracing.Tracer) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		endpoint := func(ctx context.Context, request interface{}) (response interface{}, err error) {
			q := request.(*KitfwRequest)
			method, ok := protocol.PROTOCOL_METHOD_MAP[q.Protoid]
			if ok != true {
				return nil, errors.New(fmt.Sprintf("error protoid:%d", q.Protoid))
			}

			//set log prefix "method"
			rlogger := ctx.Value("logger").(*logger.Logger)
			rlogger.With("protoid", q.Protoid, "method", method)

			v, err := s.Process(ctx, q.Protoid, q.Payload)
			if err != nil {
				return nil, err
			}
			return &KitfwReply{
				Protoid: q.Protoid + 1,
				Payload: v,
			}, nil
		}

		q := request.(*KitfwRequest)
		method, ok := protocol.PROTOCOL_METHOD_MAP[q.Protoid]
		if ok != true {
			return nil, errors.New(fmt.Sprintf("error protoid:%d", q.Protoid))
		}
		endpoint = TraceInternalService(tracer, method)(endpoint)
		return endpoint(ctx, request)
	}
}

func EndpointInstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			q := request.(*KitfwRequest)
			method, ok := protocol.PROTOCOL_METHOD_MAP[q.Protoid]
			if ok != true {
				return nil, errors.New(fmt.Sprintf("error protoid:%d", q.Protoid))
			}
			duration := duration.With("method", method, "logid", ctx.Value("logid").(string))
			endpoint := func(ctx context.Context, request interface{}) (response interface{}, err error) {
				defer func(begin time.Time) {
					duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
				}(time.Now())
				return next(ctx, request)
			}
			return endpoint(ctx, request)
		}
	}
}

func EndpointLoggingMiddleware(tracer stdopentracing.Tracer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				rlogger := ctx.Value("logger").(*logger.Logger)
				rlogger.Info("error", err, "all-took", time.Since(begin))
				rlogger = nil
			}(time.Now())
			return next(ctx, request)
		}
	}
}

func TraceInternalService(tracer stdopentracing.Tracer, operationName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			var serviceSpan stdopentracing.Span
			if parentSpan := stdopentracing.SpanFromContext(ctx); parentSpan != nil {
				serviceSpan = tracer.StartSpan(
					operationName,
					stdopentracing.ChildOf(parentSpan.Context()),
				)
			} else {
				serviceSpan = tracer.StartSpan(operationName)
			}
			defer serviceSpan.Finish()
			serviceSpan.LogKV("logid", ctx.Value("logid"))
			otext.SpanKindRPCServer.Set(serviceSpan)
			ctx = stdopentracing.ContextWithSpan(ctx, serviceSpan)
			return next(ctx, request)
		}
	}
}
