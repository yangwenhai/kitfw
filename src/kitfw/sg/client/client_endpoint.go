// Package grpc provides a gRPC client for the add service.
package main

import (
	"context"

	stdopentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"kitfw/sg/pb"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// New returns an AddService backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func NewEndPoint(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger, opname string) endpoint.Endpoint {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.
	var endpoint endpoint.Endpoint
	{
		endpoint = grpctransport.NewClient(
			conn,
			"pb.Kitfw",
			"Process",
			pb.EncodeGRPCProcessRequest,
			pb.DecodeGRPCProcessResponse,
			pb.KitfwReply{},
			grpctransport.ClientBefore(opentracing.ToGRPCRequest(tracer, logger)),
			grpctransport.ClientBefore(RequestMetaDataFunc),
		).Endpoint()
		endpoint = opentracing.TraceClient(tracer, opname)(endpoint)
	}
	return endpoint
}

func RequestMetaDataFunc(ctx context.Context, md *metadata.MD) context.Context {
	tmpmd, ok := metadata.FromContext(ctx)
	// "ToGRPCRequest" will modify the md,so we just apennd
	if ok {
		for k, v := range tmpmd {
			(*md)[k] = append((*md)[k], v...)
		}
	}
	return ctx
}
