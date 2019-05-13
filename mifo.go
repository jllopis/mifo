package mifo

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jllopis/mifo/option"
	"github.com/jllopis/mifo/server/grpc"
	"github.com/jllopis/mifo/server/grpc/middleware"
	"github.com/jllopis/mifo/server/rest"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
)

type MicroServer interface {
	Serve() error
	Shutdown()
	Version() string
	String() string
}

type Mserver struct {
	opts         *option.MsOptions
	grpcSrv      *grpc.GrpcServer
	restSrv      *rest.RestServer
	tracker      *errgroup.Group
	shutdownFunc func()
}

func NewMicroServer(options ...option.MsOption) MicroServer {
	return &Mserver{
		opts: option.New(options...),
	}
}

func (ms *Mserver) String() string {
	return fmt.Sprintf("gRPC Microservice\n=================\n  Name: %s\n  ID: %s\n", ms.opts.Name, ms.opts.ID)
}

func (ms *Mserver) Version() string {
	return ms.opts.Version
}

func (ms *Mserver) Serve() error {
	// create the tcp muxer
	mux, err := newCmux(nil, ms.opts.Address)
	if err != nil {
		log.Fatalf("cant create tcp listener for CMux, error: %s\n", err.Error())
	}
	// Match connections in order:
	// First grpc, and otherwise HTTP.
	// grpcListener := mux.Match(cmux.HTTP2HeaderFieldPrefix("content-type", "application/grpc"))
	grpcListener := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	// Otherwise, we match it againts HTTP1 methods. If matched,
	// it is sent through the "httpl" listener.
	// httpListener := mux.Match(cmux.HTTP1Fast())
	// Any significa que no hay coincidencia previa
	httpListener := mux.Match(cmux.Any())

	log.Printf("service started on %s\n", ms.opts.Address)

	// run gRPC gateway
	grpcSrv := grpc.New(ms.opts)
	go grpcSrv.Serve(grpcListener)
	ms.grpcSrv = grpcSrv

	// run HTTP gateway
	restSrv := rest.New(ms.opts)
	go restSrv.Serve(httpListener)
	ms.restSrv = restSrv

	if ms.opts.UsePrometheus {
		// Register Prometheus metrics handler.
		restSrv.GetHttpMux().Handle("/metrics", promhttp.Handler())
	}

	return mux.Serve()
}

func newCmux(l net.Listener, addr string) (cmux.CMux, error) {
	var err error
	if l == nil {
		// Create the main listener.
		l, err = net.Listen("tcp", addr)
		if err != nil {
			return nil, err
		}
	}
	return cmux.New(l), nil
}

func (ms *Mserver) Shutdown() {
	fmt.Println("Shutting server down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	ms.restSrv.Shutdown()
	middleware.HealthShutdown()

	// the last one because in our test it hangs without quittion
	ms.grpcSrv.Shutdown(shutdownCtx)
	fmt.Println("Done!")
}
