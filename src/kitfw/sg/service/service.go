package service

import (
	"errors"
	"fmt"
	protocol "kitfw/sg/protocol"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"golang.org/x/net/context"
)

type Service interface {
	Process(ctx context.Context, protoid int32, payload []byte) (context.Context, []byte, error)
}

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

type serviceLoggingMiddleware struct {
	logger log.Logger
	next   Service
}

// ServiceLoggingMiddleware returns a service middleware that logs the
// parameters and result of each method invocation.
func ServiceLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return serviceLoggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

func (mw serviceLoggingMiddleware) Process(ctx context.Context, protoid int32, payload []byte) (new_ctx context.Context, ret []byte, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"logid", new_ctx.Value("logid"),
			"protoid", protoid,
			"method", new_ctx.Value("method"),
			"userid", new_ctx.Value("userid"),
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	new_ctx, ret, err = mw.next.Process(ctx, protoid, payload)
	return new_ctx, ret, err
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

func (mw serviceInstrumentingMiddleware) Process(ctx context.Context, protoid int32, payload []byte) (new_ctx context.Context, ret []byte, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", new_ctx.Value("method").(string), "protoid", fmt.Sprint(protoid), "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	new_ctx, ret, err = mw.next.Process(ctx, protoid, payload)
	return new_ctx, ret, err
}

type basicService struct {
	logger log.Logger
}

func NewBasicService(logger log.Logger) Service {
	return basicService{logger: logger}
}

// process implements Service.
func (s basicService) Process(ctx context.Context, protoid int32, payload []byte) (context.Context, []byte, error) {
	h := GetHandler(protoid, s.logger)
	if h == nil {
		return ctx, nil, errors.New(fmt.Sprintf("error protoid:%d", protoid))
	}
	return h.Process(ctx, payload)
}

type SumHandler struct {
	request *protocol.SumRequest
	reply   *protocol.SumReply
	logger  log.Logger
}

func NewSumHandler(logger log.Logger) *SumHandler {
	request := &protocol.SumRequest{}
	reply := &protocol.SumReply{}
	return &SumHandler{request, reply, logger}
}

func (handler *SumHandler) Process(ctx context.Context, payload []byte) (context.Context, []byte, error) {
	err := protocol.Decode(handler.request, payload)
	if err != nil {
		return ctx, nil, err
	}

	ctx = context.WithValue(ctx, "method", "Sum")
	ctx = context.WithValue(ctx, "userid", handler.request.UserId)
	logger := log.NewContext(handler.logger).With("logid", ctx.Value("logid"))
	logger = log.NewContext(logger).With("logid", ctx.Value("logid"))
	logger = log.NewContext(logger).With("protoid", protocol.PROTOCOL_SUM_REQUEST)
	logger = log.NewContext(logger).With("method", "Sum")
	logger = log.NewContext(logger).With("userid", handler.request.UserId)
	handler.logger = logger

	handler.doProcess(ctx)

	ret, err := protocol.Encode(handler.reply)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, ret, nil
}

type ConcatHandler struct {
	request *protocol.ConcatRequest
	reply   *protocol.ConcatReply
	logger  log.Logger
}

func NewConcatHandler(logger log.Logger) *ConcatHandler {
	request := &protocol.ConcatRequest{}
	reply := &protocol.ConcatReply{}
	return &ConcatHandler{request, reply, logger}
}

func (handler *ConcatHandler) Process(ctx context.Context, payload []byte) (context.Context, []byte, error) {
	err := protocol.Decode(handler.request, payload)
	if err != nil {
		return ctx, nil, err
	}

	ctx = context.WithValue(ctx, "method", "Concat")
	ctx = context.WithValue(ctx, "userid", handler.request.UserId)

	handler.logger = log.NewContext(handler.logger).With("logid", ctx.Value("logid"))
	handler.logger = log.NewContext(handler.logger).With("protoid", protocol.PROTOCOL_CONCAT_REQUEST)
	handler.logger = log.NewContext(handler.logger).With("method", "Concat")
	handler.logger = log.NewContext(handler.logger).With("userid", handler.request.UserId)

	handler.doProcess(ctx)

	ret, err := protocol.Encode(handler.reply)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, ret, nil
}

type BaseHandler interface {
	Process(context.Context, []byte) (context.Context, []byte, error)
	doProcess(ctx context.Context)
}

type NewHandlerFunc func(logger log.Logger) BaseHandler

var HandlerMap = map[int32]NewHandlerFunc{
	protocol.PROTOCOL_SUM_REQUEST:    func(logger log.Logger) BaseHandler { return NewSumHandler(logger) },
	protocol.PROTOCOL_CONCAT_REQUEST: func(logger log.Logger) BaseHandler { return NewConcatHandler(logger) },
}

func GetHandler(ProtoId int32, logger log.Logger) BaseHandler {
	if HandlerMap[ProtoId] != nil {
		return HandlerMap[ProtoId](logger)
	} else {
		return nil
	}
}
