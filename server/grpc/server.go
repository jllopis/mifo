package grpc

import (
	"context"
	"errors"
	"net"
	"os"

	"github.com/jllopis/mifo/logger"
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

func (g *GrpcServer) Shutdown(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			logger.Log.Info(ctx.Err().Error()) // prints "context deadline exceeded"
			os.Exit(1)
		}
	}()
	logger.Log.Info("Shutting down gRPC server...")
	g.server.Stop()
	logger.Log.Info("gRPC server shot down")
}
