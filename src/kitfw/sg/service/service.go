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
	Process(ctx context.Context, protoid int32, payload []byte) ([]byte, error)
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

func (mw serviceLoggingMiddleware) Process(ctx context.Context, protoid int32, payload []byte) (ret []byte, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"logid", ctx.Value("logid"),
			"method", ctx.Value("method"),
			"userid", ctx.Value("userid"),
			"protoid", protoid,
			"error", err,
			"took", time.Since(begin),
		)
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
		lvs := []string{"method", "Process", "protoid", fmt.Sprint(protoid), "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.next.Process(ctx, protoid, payload)
}

type basicService struct{}

func NewBasicService() Service {
	return basicService{}
}

// process implements Service.
func (s basicService) Process(ctx context.Context, protoid int32, payload []byte) ([]byte, error) {
	h := GetHandler(protoid)
	if h == nil {
		return nil, errors.New(fmt.Sprintf("error protoid:%d", protoid))
	}
	return h.Process(ctx, payload)
}

type SumHandler struct {
	request *protocol.SumRequest
	reply   *protocol.SumReply
}

func NewSumHandler() *SumHandler {
	request := &protocol.SumRequest{}
	reply := &protocol.SumReply{}
	return &SumHandler{request, reply}
}

func (handler *SumHandler) Process(ctx context.Context, payload []byte) ([]byte, error) {
	err := protocol.Decode(handler.request, payload)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, "method", "Sum")
	ctx = context.WithValue(ctx, "userid", handler.request.UserId)

	handler.doProcess(ctx)

	ret, err := protocol.Encode(handler.reply)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type ConcatHandler struct {
	request *protocol.ConcatRequest
	reply   *protocol.ConcatReply
}

func NewConcatHandler() *ConcatHandler {
	request := &protocol.ConcatRequest{}
	reply := &protocol.ConcatReply{}
	return &ConcatHandler{request, reply}
}

func (handler *ConcatHandler) Process(ctx context.Context, payload []byte) ([]byte, error) {
	err := protocol.Decode(handler.request, payload)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, "method", "Sum")
	ctx = context.WithValue(ctx, "userid", handler.request.UserId)

	handler.doProcess(ctx)

	ret, err := protocol.Encode(handler.reply)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type BaseHandler interface {
	Process(context.Context, []byte) ([]byte, error)
	doProcess(ctx context.Context)
}

type NewHandlerFunc func() BaseHandler

var HandlerMap = map[int32]NewHandlerFunc{
	protocol.PROTOCOL_SUM_REQUEST:    func() BaseHandler { return NewSumHandler() },
	protocol.PROTOCOL_CONCAT_REQUEST: func() BaseHandler { return NewConcatHandler() },
}

func GetHandler(ProtoId int32) BaseHandler {
	if HandlerMap[ProtoId] != nil {
		return HandlerMap[ProtoId]()
	} else {
		return nil
	}
}
