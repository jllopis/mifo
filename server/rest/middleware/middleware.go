package middleware

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// A ServerOption sets options such as credentials, codec and keepalive parameters, etc.
type RestService interface {
	RegisterService(context.Context, *runtime.ServeMux, string, ...grpc.DialOption)
}

// FuncServerOption wraps a function that modifies serverOptions into an
// implementation of the ServerOption interface.
type FuncRestService func(context.Context, *runtime.ServeMux, string, ...grpc.DialOption)

func (f FuncRestService) RegisterService(ctx context.Context, r *runtime.ServeMux, addr string, opts ...grpc.DialOption) {
	if r == nil {
		panic(errors.New("option.FuncRestService.RegisterService: ServeMux server is nil"))
	}
	f(ctx, r, addr, opts...)
}

func NewFuncRestService(f func(context.Context, *runtime.ServeMux, string, ...grpc.DialOption)) FuncRestService {
	return f
}
