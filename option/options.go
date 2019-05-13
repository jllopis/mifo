package option

import (
	"fmt"
	"net/http"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jllopis/mifo/server/grpc/middleware"
	restMw "github.com/jllopis/mifo/server/rest/middleware"
	"github.com/jllopis/mifo/version"
	"github.com/justinas/alice"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type MsOptions struct {
	Name               string
	ID                 string
	Address            string
	Version            string
	Services           []middleware.GrpcService
	RestServices       []restMw.RestService
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
	GrpcOptions        []grpc.ServerOption
	RestOptions        []runtime.ServeMuxOption
	RestMiddleware     alice.Chain

	UsePrometheus bool
}

var defaultMsOptions = &MsOptions{
	ID:             ksuid.New().String(),
	Address:        ":58000",
	Version:        fmt.Sprintf("%s-%s-%s", version.Version, version.BuildDate, version.GitCommit),
	RestMiddleware: alice.New(),
}

type MsOption func(*MsOptions)

func (c *MsOptions) ServerOptions() []grpc.ServerOption {
	opts := append(
		[]grpc.ServerOption{
			grpc_middleware.WithUnaryServerChain(c.UnaryInterceptors...),
			grpc_middleware.WithStreamServerChain(c.StreamInterceptors...),
		},
		c.GrpcOptions...,
	)
	return opts
}

func New(options ...MsOption) *MsOptions {
	mso := defaultMsOptions
	for _, f := range options {
		f(mso)
	}
	return mso
}

func Name(name string) MsOption {
	return func(o *MsOptions) {
		o.Name = name
	}
}

func ID(id string) MsOption {
	return func(o *MsOptions) {
		o.ID = id
	}
}

func Address(address string) MsOption {
	return func(o *MsOptions) {
		o.Address = address
	}
}

// WithStatsHandler ConnectionTimeout returns a ServerOption that sets the timeout for connection establishment (up to and including HTTP/2 handshaking) for all new connections.
// If this is not set, the default is 120 seconds.
func WithConnTimeout(t time.Duration) MsOption {
	return func(c *MsOptions) {
		c.GrpcOptions = append(c.GrpcOptions, grpc.ConnectionTimeout(t))
	}
}

// WithMaxConcurrentStreams returns a ServerOption that will apply a limit on the number
// of concurrent streams to each ServerTransport.
func WithMaxConcurrentStreams(num uint32) MsOption {
	return func(c *MsOptions) {
		c.GrpcOptions = append(c.GrpcOptions, grpc.MaxConcurrentStreams(num))
	}
}

// WithReflection adds grpc server reflection to the list of Services ref: https://godoc.org/google.golang.org/grpc/reflection
func WithReflection() MsOption {
	return func(o *MsOptions) {
		o.Services = append(o.Services, middleware.NewReflection())
	}
}

func WithHealthz() MsOption {
	return func(o *MsOptions) {
		o.Services = append(o.Services, middleware.NewHealthz())
	}
}

func WithServices(services ...middleware.GrpcService) MsOption {
	return func(o *MsOptions) {
		o.Services = append(o.Services, services...)
	}
}

func WithRestServices(services ...restMw.RestService) MsOption {
	return func(o *MsOptions) {
		o.RestServices = append(o.RestServices, services...)
	}
}

func WitUnaryLoggingMiddleware(logger *zap.Logger) MsOption {
	return func(o *MsOptions) {
		o.Services = append(o.Services, middleware.NewLogger(logger))
		o.UnaryInterceptors = append(o.UnaryInterceptors, middleware.NewUnaryLogger(logger))
	}
}

func WitStreamLoggingMiddleware(logger *zap.Logger) MsOption {
	return func(o *MsOptions) {
		o.Services = append(o.Services, middleware.NewLogger(logger))
		o.StreamInterceptors = append(o.StreamInterceptors, middleware.NewStreamLogger(logger))
	}
}

func WitUnaryPrometheusMiddleware() MsOption {
	return func(o *MsOptions) {
		o.Services = append(o.Services, middleware.NewPrometheus())
		o.UnaryInterceptors = append(o.UnaryInterceptors, middleware.NewUnaryPrometheus())
		o.UsePrometheus = true
	}
}

func WitStreamPrometheusMiddleware() MsOption {
	return func(o *MsOptions) {
		o.Services = append(o.Services, middleware.NewPrometheus())
		o.StreamInterceptors = append(o.StreamInterceptors, middleware.NewStreamPrometheus())
		o.UsePrometheus = true
	}
}

// grpc-gateway options
func WithRestMarshaler(mime string, msh runtime.Marshaler) MsOption {
	return func(o *MsOptions) {
		if mime == "" {
			mime = runtime.MIMEWildcard
		}
		if msh == nil {
			return
		}
		o.RestOptions = append(o.RestOptions, runtime.WithMarshalerOption(mime, msh))
	}
}

func WithRestIncomingHeaderMatcher(matcher runtime.HeaderMatcherFunc) MsOption {
	return func(o *MsOptions) {
		if matcher == nil {
			matcher = runtime.DefaultHeaderMatcher
		}
		o.RestOptions = append(o.RestOptions, runtime.WithIncomingHeaderMatcher(matcher))
	}
}

func WithRestOutgoingHeaderMatcher(matcher runtime.HeaderMatcherFunc) MsOption {
	return func(o *MsOptions) {
		if matcher == nil {
			matcher = runtime.DefaultHeaderMatcher
		}
		o.RestOptions = append(o.RestOptions, runtime.WithOutgoingHeaderMatcher(matcher))
	}
}

func WithRestLoggingMiddleware(logger *zap.Logger) MsOption {
	return func(o *MsOptions) {
		o.RestMiddleware = o.RestMiddleware.Append(restMw.AddRequestID(), restMw.AddLogger(logger))
	}
}

func WithHTTPMiddleware(mw func(http.Handler) http.Handler) MsOption {
	return func(o *MsOptions) {
		o.RestMiddleware = o.RestMiddleware.Append(mw)
	}
}
