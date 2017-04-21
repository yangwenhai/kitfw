package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"context"

	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"

	logger "kitfw/sg/log"
	"kitfw/sg/pb"
	kitservice "kitfw/sg/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	sdetcd "github.com/go-kit/kit/sd/etcd"
	"github.com/go-kit/kit/tracing/opentracing"
)

func main() {
	var (
		debugAddr  = flag.String("debugAddr", ":8080", "Debug and metrics listen address")
		grpcAddr   = flag.String("grpcAddr", "127.0.0.1:8081", "gRPC (HTTP) listen address")
		zipkinAddr = flag.String("zipkinAddr", "", "Enable Zipkin tracing via a Kafka server host:port")
		etcdAddr   = flag.String("etcdAddr", "", "gRPC (HTTP) listen address")
	)
	flag.Parse()

	//log
	logger.SetDefaultLogLevel(logger.LevelDebug)
	logger.Info("msg", "hello kitfw")
	defer logger.Info("msg", "goodbye kitfw")

	// Metrics domain.
	fieldKeys := []string{"method", "protoid", "error"}
	var requestCount metrics.Counter
	{
		// Business level metrics.
		requestCount = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "kitfw",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys)
	}
	var duration metrics.Histogram
	{
		// Transport level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "kitfw",
			Name:      "request_duration_ns",
			Help:      "Request duration in nanoseconds.",
		}, fieldKeys)
	}
	var endpointDuration metrics.Histogram
	{
		// Transport level metrics.
		endpointDuration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "kitfw",
			Name:      "endpoint_request_duration_ns",
			Help:      "endpoint request duration in nanoseconds.",
		}, []string{"method", "success"})
	}

	// Tracing domain.
	var tracer stdopentracing.Tracer
	{
		if *zipkinAddr != "" {
			logger.Info("tracer", "Zipkin", "zipkinAddr", *zipkinAddr)
			// collector, err := zipkin.NewKafkaCollector(
			// 	strings.Split(*zipkinAddr, ","),
			// 	zipkin.KafkaLogger(logger),
			// )
			collector, err := zipkin.NewHTTPCollector(*zipkinAddr)
			if err != nil {
				logger.Error("tracer", "Zipkin", "err", err)
				os.Exit(1)
			}
			tracer, err = zipkin.NewTracer(
				zipkin.NewRecorder(collector, true, "servicename:host_ip", "kitfw"), zipkin.WithLogger(logger.GetDefaultLogger()),
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

	// Business domain.
	var service kitservice.Service
	{
		service = kitservice.NewBasicService()
		service = kitservice.ServiceLoggingMiddleware()(service)
		service = kitservice.ServiceInstrumentingMiddleware(requestCount, duration)(service)
	}

	// Endpoint domain.
	var requestEndpoint endpoint.Endpoint
	{
		processDuration := endpointDuration.With("method", "Process")

		requestEndpoint = pb.MakeProcessEndpoint(service)
		requestEndpoint = opentracing.TraceServer(tracer, "Process")(requestEndpoint)
		requestEndpoint = pb.EndpointInstrumentingMiddleware(processDuration)(requestEndpoint)
		requestEndpoint = pb.EndpointLoggingMiddleware()(requestEndpoint)
	}

	// Mechanical domain.
	errc := make(chan error)
	ctx := context.Background()

	// Interrupt handler.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// Debug listener.
	go func() {
		m := http.NewServeMux()
		m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		m.Handle("/metrics", stdprometheus.Handler())

		logger.Info("transport", "debug", "debugAddr", *debugAddr)
		errc <- http.ListenAndServe(*debugAddr, m)
	}()

	// gRPC transport.
	go func() {
		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}

		srv := pb.MakeGRPCServer(ctx, requestEndpoint, tracer, logger.GetDefaultLogger())
		s := grpc.NewServer()
		pb.RegisterKitfwServer(s, srv)

		logger.Info("transport", "gRPC", "grpcAddr", *grpcAddr)
		errc <- s.Serve(ln)
	}()

	// etcd registry
	var etcdClient sdetcd.Client
	if *etcdAddr != "" {
		var etcdErr error
		etcdClient, etcdErr = sdetcd.NewClient(
			context.Background(),
			[]string{*etcdAddr},
			sdetcd.ClientOptions{
				DialTimeout:             2 * time.Second,
				DialKeepAlive:           2 * time.Second,
				HeaderTimeoutPerRequest: 2 * time.Second,
			},
		)
		if etcdErr != nil {
			logger.Error("unexpected error creating client", etcdErr)
			os.Exit(1)
		}
		if etcdClient == nil {
			logger.Error("expected new Client, got nil")
			os.Exit(1)
		}
		// registrar
		go func() {
			key := fmt.Sprintf("%s/%s", "/kitfw/service", *grpcAddr)
			registrar := sdetcd.NewRegistrar(etcdClient, sdetcd.Service{
				Key:           key,
				Value:         *grpcAddr,
				DeleteOptions: nil,
			}, logger.GetDefaultLogger())

			if registrar == nil {
				logger.Error("expected new Client, got nil")
				os.Exit(1)
			}
			registrar.Deregister()
			registrar.Register()
			for {
				time.Sleep(8 * time.Second)
				res, _ := etcdClient.GetEntries(key)
				if res == nil {
					registrar.Register()
				}
			}
		}()
		logger.Info("etcdAddr", *etcdAddr)
	}

	// Run!
	logger.Error("exit", <-errc)
}
