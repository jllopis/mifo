package middleware

import (
	"errors"

	"google.golang.org/grpc"
)

// A ServerOption sets options such as credentials, codec and keepalive parameters, etc.
type GrpcService interface {
	RegisterService(*grpc.Server)
}

// EmptyServerOption does not alter the server configuration. It can be embedded
// in another structure to build custom server options.
type EmptyGrpcService struct{}

func (EmptyGrpcService) RegisterService(*grpc.Server) {}

// FuncServerOption wraps a function that modifies serverOptions into an
// implementation of the ServerOption interface.
type FuncGrpcService func(*grpc.Server)

func (f FuncGrpcService) RegisterService(s *grpc.Server) {
	if s == nil {
		panic(errors.New("option.RegisterService: grpc server is nil"))
	}
	f(s)
}

func NewFuncGrpcService(f func(*grpc.Server)) FuncGrpcService {
	return f
}
