package service

import (
	"errors"
	"fmt"
	logger "kitfw/sg/log"
	protocol "kitfw/sg/protocol"
	"time"

	"context"

	"github.com/go-kit/kit/metrics"
)

type Service interface {
	Process(ctx context.Context, protoid int32, payload []byte) ([]byte, error)
}

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

type serviceLoggingMiddleware struct {
	next Service
}

// ServiceLoggingMiddleware returns a service middleware that logs the
// parameters and result of each method invocation.
func ServiceLoggingMiddleware() Middleware {
	return func(next Service) Service {
		return serviceLoggingMiddleware{
			next: next,
		}
	}
}

func (mw serviceLoggingMiddleware) Process(ctx context.Context, protoid int32, payload []byte) (ret []byte, err error) {
	defer func(begin time.Time) {
		rlogger := ctx.Value("logger").(*logger.Logger)
		rlogger.Info("error", err, "service-took", time.Since(begin))
	}(time.Now())
	return mw.next.Process(ctx, protoid, payload)
}

type serviceInstrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

// ServiceInstrumentingMiddleware returns a service middleware that instruments
// the number of integers summed and characters concatenated over the lifetime of
// the service.
func ServiceInstrumentingMiddleware(rc metrics.Counter, rl metrics.Histogram) Middleware {
	return func(next Service) Service {
		return serviceInstrumentingMiddleware{
			requestCount:   rc,
			requestLatency: rl,
			next:           next,
		}
	}
}

func (mw serviceInstrumentingMiddleware) Process(ctx context.Context, protoid int32, payload []byte) (ret []byte, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", protocol.PROTOCOL_METHOD_MAP[protoid], "protoid", fmt.Sprint(protoid), "logid", ctx.Value("logid").(string), "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.next.Process(ctx, protoid, payload)
}

type basicService struct {
}

func NewBasicService() Service {
	return basicService{}
}

// process implements Service.
func (s basicService) Process(ctx context.Context, protoid int32, payload []byte) ([]byte, error) {
	h := GetHandler(protoid)
	if h == nil {
		return nil, errors.New(fmt.Sprintf("invalid protoid:%d", protoid))
	}
	return h.Process(ctx, payload)
}
