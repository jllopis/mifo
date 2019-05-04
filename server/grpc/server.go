package grpc

import (
	"errors"
	"net"

	"github.com/jllopis/mifo/option"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	server *grpc.Server
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
func (g *GrpcServer) Serve(l net.Listener) {
	g.server.Serve(l)
}
