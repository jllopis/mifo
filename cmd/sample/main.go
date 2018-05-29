package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/acbapis/acbapis/status"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jllopis/mifo"
	"github.com/jllopis/mifo/cmd/sample/impl"
	"github.com/jllopis/mifo/log"
	"github.com/jllopis/mifo/version"
)

func main() {
	log.Info("Starting sample server")
	log.Info("Server version", "value", version.SemVer.String())
	log.Info("Git tag", "value", version.GitCommit)
	log.Info("Build", "date", version.BuildDate)

	// var err error
	// var l net.Listener
	// if l == nil {
	// 	// Create the main listener.
	// 	l, err = net.Listen("tcp", ":3000")
	// 	if err != nil {
	// 		log.Err("cant create CMux", "error", err.Error())
	// 		os.Exit(2)
	// 	}
	// }

	// api := mifo.NewMicroServer().WithListener(l)
	api := mifo.NewMicroServer().SetPort(58000)
	if api == nil {
		os.Exit(1)
	}

	log.Info("Server info", "name", api.Name, "port", api.Port)
	fmt.Println(api.String())

	// Register middleware interceptors
	api.UsePrometheus(&mifo.PrometheusInterceptorConfig{EnableTimeHistogram: true})
	api.UseReflection()

	// Register gRPC services
	statusRegistrator := func(srv *grpc.Server) {
		status.RegisterStatusServiceServer(srv, &impl.StatusService{})
	}
	api.Register(statusRegistrator)

	// Use grpc-gw ...
	api.UseGrpcGw()
	// ... and register http proxy services
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	statusGwRegistrator := func(srv *runtime.ServeMux) {
		status.RegisterStatusServiceHandlerFromEndpoint(context.Background(), srv, "localhost:"+api.Port, opts)
	}
	api.RegisterGw(statusGwRegistrator)

	log.Info("Serving...")
	log.Info("Prometheus metrics on port " + api.Port)
	go api.Serve()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sc
	api.Stop()

	log.Info("Sample finished")
}
