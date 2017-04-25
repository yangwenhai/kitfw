package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"kitfw/sg/pb"

	protocol "kitfw/sg/protocol"

	"context"

	"kitfw/sg/define"
	logger "kitfw/sg/log"

	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	address = "localhost:8081"
)

func main() {

	var (
		grpcAddr   = flag.String("grpcAddr", "127.0.0.1:8081", "gRPC (HTTP) listen address")
		zipkinAddr = flag.String("zipkinAddr", "", "Enable Zipkin tracing via a Kafka server host:port")
	)
	flag.Parse()

	logger.SetDefaultLogLevel(logger.LevelDebug)

	// Set up a connection to the server.
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKitfwClient(conn)
	_ = c

	// Tracing domain.
	var tracer stdopentracing.Tracer
	{
		if *zipkinAddr != "" {
			logger.Info("tracer", "Zipkin", "zipkinAddr", *zipkinAddr)
			// collector, err := zipkin.NewKafkaCollector(
			// 	strings.Split(*zipkinAddr, ","),
			// 	zipkin.KafkaLogger(logger),
			// )
			collector, err := zipkin.NewHTTPCollector(*zipkinAddr, zipkin.HTTPBatchSize(1))
			if err != nil {
				logger.Error("tracer", "Zipkin", "err", err)
				os.Exit(1)
			}
			tracer, err = zipkin.NewTracer(
				zipkin.NewRecorder(collector, false, "0.0.0.0:8081", define.SERVER_NAME),
				zipkin.ClientServerSameSpan(true),
			)
			if err != nil {
				logger.Error("tracer", "Zipkin", "err", err)
				os.Exit(1)
			}
		} else {
			logger.Info("tracer", "none")
			tracer = stdopentracing.GlobalTracer() // no-op
		}
	}

	sum(conn, tracer, 1000, 2000)

	concat(conn, tracer, "hello", "world")

	//waitting for tracer finish send msg to zipkin
	time.Sleep(5 * time.Second)

}

func concat(conn *grpc.ClientConn, tracer stdopentracing.Tracer, str1 string, str2 string) {

	rand.Seed(time.Now().UnixNano())
	userid := int64(rand.Intn(10000))%10000 + 10000

	//encode capnp
	req := &protocol.ConcatRequest{
		UserId: userid,
		Str1:   str1,
		Str2:   str2,
	}
	payload, _ := protocol.Encode(req)

	//send request
	starttime := time.Now()
	logid := fmt.Sprintf("%d", time.Now().UnixNano()%10000000)
	md := metadata.New(map[string]string{"userid": fmt.Sprintf("%d", userid), "logid": logid})
	// create a new context with this metadata
	ctx := metadata.NewContext(context.Background(), md)

	//endpoint
	sumendpoint := NewEndPoint(conn, tracer, logger.GetDefaultLogger(), "Concat")
	r, err := sumendpoint(ctx, &pb.KitfwRequest{
		Protoid: protocol.PROTOCOL_CONCAT_REQUEST,
		Payload: payload,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	//decode capnp
	res := &protocol.ConcatReply{}
	err = protocol.Decode(res, r.(*pb.KitfwReply).Payload)
	if err != nil {
		fmt.Println("err!", err)
		return
	}

	log.Printf("concat response: %s took: %fms", res.Val, time.Since(starttime).Seconds()*1000)
}

func sum(conn *grpc.ClientConn, tracer stdopentracing.Tracer, num1 int64, num2 int64) {

	rand.Seed(time.Now().UnixNano())
	userid := int64(rand.Intn(10000))%10000 + 10000

	//encode capnp
	req := &protocol.SumRequest{
		UserId: userid,
		Num1:   num1,
		Num2:   num2,
	}
	payload, _ := protocol.Encode(req)

	//send request
	starttime := time.Now()
	logid := fmt.Sprintf("%d", time.Now().UnixNano()%10000000)
	md := metadata.New(map[string]string{"userid": fmt.Sprintf("%d", userid), "logid": logid})
	// create a new context with this metadata
	ctx := metadata.NewContext(context.Background(), md)

	//endpoint
	sumendpoint := NewEndPoint(conn, tracer, logger.GetDefaultLogger(), "Sum")
	r, err := sumendpoint(ctx, &pb.KitfwRequest{
		Protoid: protocol.PROTOCOL_SUM_REQUEST,
		Payload: payload,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	//decode capnp
	res := &protocol.SumReply{}
	err = protocol.Decode(res, r.(*pb.KitfwReply).Payload)
	if err != nil {
		fmt.Println("err!", err)
		return
	}

	log.Printf("sum response: %d took: %fms", res.Val, time.Since(starttime).Seconds()*1000)
}
