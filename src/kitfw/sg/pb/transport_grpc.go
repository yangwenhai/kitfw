package pb

// This file provides server-side bindings for the gRPC transport.
// It utilizes the transport/grpc.Server.

import (
	stdopentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC AddServer.
func MakeGRPCServer(ctx context.Context, endpoint endpoint.Endpoint, tracer stdopentracing.Tracer, logger log.Logger) KitfwServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		process: grpctransport.NewServer(
			ctx,
			endpoint,
			DecodeGRPCProcessRequest,
			EncodeGRPCProcessResponse,
			append(options, grpctransport.ServerBefore(opentracing.FromGRPCRequest(tracer, "Proces", logger)))...,
		),
	}
}

type grpcServer struct {
	process grpctransport.Handler
}

func (s *grpcServer) Process(ctx context.Context, req *KitfwRequest) (*KitfwReply, error) {
	_, res, err := s.process.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*KitfwReply), nil
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
