package pb

// This file provides server-side bindings for the gRPC transport.
// It utilizes the transport/grpc.Server.

import (
	"context"

	"google.golang.org/grpc/metadata"

	logger "kitfw/sg/log"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	oldcontext "golang.org/x/net/context"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC AddServer.
func MakeGRPCServer(ctx context.Context, endpoint endpoint.Endpoint, tracer stdopentracing.Tracer, logger log.Logger) KitfwServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		process: grpctransport.NewServer(
			endpoint,
			DecodeGRPCProcessRequest,
			EncodeGRPCProcessResponse,
			append(options, grpctransport.ServerBefore(RequestMetaDataFunc), grpctransport.ServerBefore(opentracing.FromGRPCRequest(tracer, "Process", logger)))...,
		),
	}
}

type grpcServer struct {
	process grpctransport.Handler
}

func (s *grpcServer) Process(ctx oldcontext.Context, req *KitfwRequest) (*KitfwReply, error) {
	_, res, err := s.process.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*KitfwReply), nil
}

func RequestMetaDataFunc(ctx context.Context, md metadata.MD) context.Context {
	logid := ""
	if logids, ok := md["logid"]; ok && len(logids) > 0 {
		logid = logids[0]
	}
	userid := ""
	if userids, ok := md["userid"]; ok && len(userids) > 0 {
		userid = userids[0]
	}
	logger.Info("logid", logid, "userid", userid)
	rlogger := logger.NewLogger()
	rlogger.SetLogLevel(logger.LevelDebug)
	rlogger.With("logid", logid, "userid", userid)
	ctx = context.WithValue(ctx, "logger", rlogger)
	ctx = context.WithValue(ctx, "logid", logid)
	return ctx
}

func DecodeGRPCProcessRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq, nil
}

func DecodeGRPCProcessResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	return grpcReply, nil
}

func EncodeGRPCProcessResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response, nil
}

func EncodeGRPCProcessRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}
