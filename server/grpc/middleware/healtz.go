package middleware

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Healthz struct {
	FuncGrpcService
}

func NewHealthz() *Healthz {
	return &Healthz{
		FuncGrpcService: func(s *grpc.Server) {
			grpc_health_v1.RegisterHealthServer(s, health.NewServer())
		},
	}
}
