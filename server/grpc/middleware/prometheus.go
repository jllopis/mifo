package middleware

import (
	"sync"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

var once sync.Once

type Prometheus struct {
	FuncGrpcService
}

func NewPrometheus() *Prometheus {
	return &Prometheus{
		FuncGrpcService: func(s *grpc.Server) {
			once.Do(func() {
				// grpc_prometheus.EnableHandlingTimeHistogram()
				grpc_prometheus.Register(s)
			})
		},
	}
}

func NewUnaryPrometheus() grpc.UnaryServerInterceptor {
	return grpc_prometheus.UnaryServerInterceptor
}

func NewStreamPrometheus() grpc.StreamServerInterceptor {
	return grpc_prometheus.StreamServerInterceptor
}
