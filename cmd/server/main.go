package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/jllopis/mifo"
	"github.com/jllopis/mifo/cmd/server/impl"
	"github.com/jllopis/mifo/logger"
	"github.com/jllopis/mifo/option"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

// Config is configuration for Server
type Config struct {
	// Service default port shared by protocols gRPC, http
	DefaultPort string `getconf:"port, default: 9090, info: port to bind services"`
	// Log parameters section
	// LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	LogLevel int `getconf:"log-level, default: 0, info: Global log level"`
	// LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
	LogTimeFormat string `getconf:"log-time-format, info: Print time format for logger e.g. 2006-01-02T15:04:05Z07:00"`
	// set max concurrent streams served by gRPC server
	GrpcMaxConcurrentStreams int `getconf:"grpc-max-concurrent-streams, default: 250, info: grpc server option MaxConcurrentStreams"`
	// Expose Prometheus metrics
	UsePrometheus bool `getconf:"use-prometheus, default: true, info: expose prometheus metrics"`
	// endpoint for Prometheus metrics
	PrometheusEndpoint string `getconf:"prometheus-endpoint, default: /metrics, info: endpoint to serve prometheus metrics. Only available if UsePrometheus()"`
	// if mode is dev, log everything
	Mode string `getconf:"mode, default: dev, info: set the environment (dev staging or prod)"`
}

func main() {
	// initialize logger
	if err := logger.Init(-1, "2006-01-02T15:04:05.999999999Z07:00"); err != nil {
		fmt.Errorf("failed to initialize logger: %v", err)
	}

	logger.Log.Info("Starting sample server")

	// api := mifo.NewMicroServer().WithListener(l)
	srv := mifo.NewMicroServer(
		// mifo server options
		option.Name("TheClass"),
		option.Address(":58000"),

		// grpc server options
		option.WithReflection(),
		option.WithHealthz(),

		// grpc unary interceptors
		option.WitUnaryLoggingMiddleware(logger.Log),

		// grpc stream interceptors
		option.WitStreamLoggingMiddleware(logger.Log),

		// grpc services registration
		option.WithServices(impl.NewStatusService()),

		// after other services, install prometheus middleware
		option.WitUnaryPrometheusMiddleware(),
		option.WitStreamPrometheusMiddleware(),
		// grpc-gw services registration
		option.WithRestServices(impl.NewStatusRestService()),

		// grpc-gw REST middleware
		option.WithHTTPMiddleware(cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
			AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
			AllowCredentials: true,
			// Enable Debugging for testing, consider disabling in production
			Debug: true,
		}).Handler),
		option.WithRestLoggingMiddleware(logger.Log),
	)
	if srv == nil {
		os.Exit(1)
	}

	logger.Log.Info("Server version", zap.String("value", srv.Version()))
	fmt.Println(srv.String())

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		s := <-sc
		fmt.Printf("[main] Signal %+v captured. Calling srv.Shutdown()\n", s)
		srv.Shutdown()
		os.Exit(0)
	}()

	logger.Log.Info(srv.Serve().Error())

	logger.Log.Info("Sample finished")
}
