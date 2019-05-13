package middleware

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var healthSrv *health.Server

type Healthz struct {
	FuncGrpcService
}

func NewHealthz() *Healthz {
	healthSrv = health.NewServer()
	return &Healthz{
		FuncGrpcService: func(s *grpc.Server) {
			grpc_health_v1.RegisterHealthServer(s, healthSrv)
		},
	}
}

func SetHealthStatus(serviceName string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	fmt.Printf("called SetHealthStatus. sn=%s, status=%d\n", serviceName, status)
	healthSrv.SetServingStatus(fmt.Sprintf("grpc.health.v1.%s", serviceName), status)
}

func HealthShutdown() {
	if healthSrv != nil {
		log.Printf("Shutting down healthz server...")
		healthSrv.Shutdown()
		log.Print("healthz server shot down")
	}
}
