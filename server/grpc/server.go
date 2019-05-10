package grpc

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/jllopis/mifo/option"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	server   *grpc.Server
	listener net.Listener
}

func New(o *option.MsOptions) *GrpcServer {
	s := grpc.NewServer(o.ServerOptions()...)
	for _, service := range o.Services {
		service.RegisterService(s)
	}

	gserver := &GrpcServer{
		server: s,
	}
	if gserver == nil {
		panic(errors.New("failed to validate grpc server"))
	}
	return gserver
}

// Serve runs gRPC service
func (g *GrpcServer) Serve(l net.Listener) error {
	g.listener = l
	return g.server.Serve(l)
}

func (g *GrpcServer) Shutdown() {
	stopped := make(chan struct{})
	go func() {
		fmt.Println("Shutting down gRPC server...")
		g.server.GracefulStop()
		close(stopped)
	}()
	fmt.Println("waiting for gRPC connections to finish...")
	t := time.NewTimer(10 * time.Second)
	select {
	case <-t.C:
		fmt.Println("timed out. stopping gRPC server...")
		g.server.Stop()
	case <-stopped:
		fmt.Println("cancelling timer...")
		t.Stop()
	}
	fmt.Println("gRPC server stopped")
	// fmt.Println("Closing grpc listener")
	// g.listener.Close()
	fmt.Println("gRPC server shot down")
}
